package main

import "github.com/dragonator/grpc-golang/internal/pb"

var employees = []*pb.Employee{
	{
		Id:          1,
		BadgeNumber: 1234,
		FirstName:   "Angela",
		LastName:    "Robins",
	},
	{
		Id:          2,
		BadgeNumber: 2345,
		FirstName:   "Yasmin",
		LastName:    "Wise",
	},
	{
		Id:          3,
		BadgeNumber: 3456,
		FirstName:   "Renee",
		LastName:    "Rees",
	},
	{
		Id:          4,
		BadgeNumber: 4567,
		FirstName:   "Amara",
		LastName:    "Becker",
	},
	{
		Id:          5,
		BadgeNumber: 5678,
		FirstName:   "Jasleen",
		LastName:    "Hope",
	},
}
