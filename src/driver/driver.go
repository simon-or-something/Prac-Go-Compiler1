package driver

type Driver struct {
	Infile string
}

func RunDriver(args []string) (*Driver, error) {
	argparser := Driver{Infile: args[1]}
	return &argparser, nil
}
