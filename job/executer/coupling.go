package executer

var executer *Executer

func GetInstance() *Executer {
	if executer == nil {
		executer = &Executer{}
	}

	return executer
}
