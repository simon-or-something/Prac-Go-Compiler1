package driver

type driver struct {
		infile string
}

func RunDriver(args []string) (*driver, error) {
		argparser := driver{infile: args[1]}
		return &argparser, nil
}


