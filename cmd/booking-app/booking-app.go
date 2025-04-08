package main

import (
	"fmt"
	"regexp"
	"sync"
	"time"
)

const errUserInput string = "Error reading user input"

var waitGroup = sync.WaitGroup{}

func main() {
	const ticketsTotal int8 = 50
	var bookings = make([]Booking, 0, 16)
	var ticketsAvailable = ticketsTotal

	for ticketsAvailable > 0 {
		var email string
		var name string
		var tickets int8

		fmt.Println("Welcome to \"Conference 9\" booking application")
		fmt.Println("We have", ticketsTotal, "tickets, and", ticketsAvailable, "tickets available.")
		name = getName()
		email = getEmail()
		tickets = getTickets()

		if tickets <= ticketsAvailable {
			var booking = Booking{
				name:    name,
				email:   email,
				tickets: tickets,
			}

			bookings = append(bookings, booking)
			ticketsAvailable -= tickets

			fmt.Println(
				"Thank you",
				name,
				"for b",
				tickets,
				"tickets. You will receive a confirmation email at",
				email,
			)
			fmt.Println("Current bookings:")

			for i, b := range bookings {
				fmt.Printf("%v: %v\n", i+1, b.name)
			}

			waitGroup.Add(1)
			go sendTicket(booking)
		} else {
			fmt.Printf(
				"There are not enough tickets remaining (%v) to book %v tickets.\n",
				ticketsAvailable,
				tickets,
			)
		}
	}

	fmt.Println("The conference is sold out. Thank you!")
	waitGroup.Wait()
}

func getEmail() string {
	for {
		var userEmail string

		fmt.Print("Enter your email address: ")
		if _, err := fmt.Scan(&userEmail); err != nil {
			fmt.Println(errUserInput, err)
			continue
		}

		if match, _ := regexp.MatchString("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", userEmail); !match {
			fmt.Printf("Invalid email address: \"%v\"\n", userEmail)
		} else {
			return userEmail
		}
	}
}

func getName() string {
	const pattern string = "^[a-zA-Z0-9.-_ ]+$"

	for {
		var userName string

		fmt.Print("Enter your preferred name: ")
		if _, err := fmt.Scan(&userName); err != nil {
			fmt.Println(errUserInput, err)
			continue
		}

		if match, _ := regexp.MatchString(pattern, userName); !match {
			fmt.Printf("Name must match the regular expression: \"%v\"\n", pattern)
		} else {
			return userName
		}
	}
}

func getTickets() int8 {
	for {
		var userTickets int8

		fmt.Print("Enter the number of tickets you want to book: ")
		if _, errUserTickets := fmt.Scan(&userTickets); errUserTickets != nil {
			fmt.Println(errUserInput, errUserTickets)
			continue
		}

		if userTickets < 1 {
			fmt.Printf("Invalid number of tickets: %d\n", userTickets)
		} else {
			return userTickets
		}
	}
}

func sendTicket(booking Booking) {
	var ticket = fmt.Sprintf("%v ticket(s) for %v", booking.tickets, booking.name)

	time.Sleep(5 * time.Second)

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Sending ticket:", ticket)
	fmt.Println("--------------------------------------------------------------------------------")

	waitGroup.Done()
}

type Booking struct {
	name    string
	email   string
	tickets int8
}
