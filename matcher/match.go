package matcher

import (
	"errors"
	"math/rand"
	"time"

	"github.com/nathan-osman/go-pearup/db"
)

func (m *Matcher) match(conn *db.Conn, p *db.Pearup) error {
	m.log.Infof("begining %s...", p.Name)

	// Seed the random number generator
	rand.Seed(time.Now().Unix())

	// Set the pear-up to complete
	p.IsComplete = true
	if err := conn.Save(p).Error; err != nil {
		return err
	}

	// Load the list of previous matches so they can be avoided
	previousMatches := []*db.Match{}
	if err := conn.
		Find(&previousMatches).
		Error; err != nil {
		return err
	}

	// Function for validating matches
	isValidMatch := func(user1ID, user2ID int64) bool {
		for _, m := range previousMatches {
			if user1ID == m.User1ID && user2ID == m.User2ID ||
				user1ID == m.User2ID && user2ID == m.User1ID {
				return false
			}
		}
		return true
	}

	// Fetch all registrations (and preload the user structs for each)
	registrations := []*db.Registration{}
	if err := conn.
		Where("pearup_id = ?", p.ID).
		Order("date").
		Preload("User").
		Find(&registrations).Error; err != nil {
		return err
	}

	// Sort the men and women into separate lists based on gender
	var (
		men   = []*db.User{}
		women = []*db.User{}
	)
	for _, r := range registrations {
		if r.User.IsMale {
			men = append(men, r.User)
		}
		if r.User.IsFemale {
			women = append(women, r.User)
		}
	}
	m.log.Infof("pairing up %d men & %d women...", len(men), len(women))

	// We've run into a rather unfortunate problem if either list is empty
	if len(men) == 0 || len(women) == 0 {
		m.log.Warn("the list for one gender is empty - no matches")
		return nil
	}

	// Determine which list has the most items
	var (
		larger  []*db.User
		smaller []*db.User
	)
	if len(men) >= len(women) {
		larger = men
		smaller = women
	} else {
		larger = women
		smaller = men
	}

	// Do this a few times until it works or give up
outer:
	for n := 0; n < 10; n++ {

		m.log.Infof("attempt #%d to pair...", n+1)

		matches := []*db.Match{}

		// Shuffle the lists (this is the random part)
		m.log.Info("shuffling lists of users...")
		for _, l := range [][]*db.User{larger, smaller} {
			for i := range l {
				j := rand.Intn(i + 1)
				l[i], l[j] = l[j], l[i]
			}
		}

		// Step through the users in the larger list and pair them with the
		// smaller list until complete
		i := 0
		for _, u := range larger {

			// Move on and try again if this match is invalid
			if !isValidMatch(u.ID, smaller[i].ID) {
				m.log.Infof(
					"%s & %s have already been paired",
					u.FacebookName,
					smaller[i].FacebookName,
				)
				continue outer
			}

			// Append the match
			matches = append(matches, &db.Match{
				PearupID: p.ID,
				User1ID:  u.ID,
				User2ID:  smaller[i].ID,
			})

			i++
			if i == len(smaller) {
				i = 0
			}
		}

		// Save all the matches
		for _, m := range matches {
			if err := conn.Save(m).Error; err != nil {
				return err
			}
		}

		m.log.Infof("completed %s", p.Name)
		return nil
	}

	return errors.New("failed to find a valid solution")
}
