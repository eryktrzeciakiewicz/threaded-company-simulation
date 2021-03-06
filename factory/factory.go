package factory

import (
	"math/rand"

	"github.com/projects/threaded-company-simulation/agents"
	"github.com/projects/threaded-company-simulation/config"
)

func RunLists() {
	go agents.SynchronizeWarehouse()
	go agents.SynchronizeTaskList()
}

func RunLogger() {
	go agents.RunLogger()
}

func RunBoss() {
	Boss := agents.Boss{1, agents.TaskListWrite, agents.LogChannel}
	go Boss.Run()
}

func CreateMultMachines() []*agents.MultiplicationMachine {
	var machines = make([]*agents.MultiplicationMachine, 0)
	for i := 0; i < config.NUM_MULT_MACHINES; i++ {
		machine := agents.MultiplicationMachine{
			Id:              i,
			Input:           make(chan agents.MachineWriteOp),
			Logger:          agents.LogChannel,
			IsBroken:        false,
			BreakdownNumber: 0,
			FixMe:           make(chan bool)}

		go machine.RunMultiplicationMachine()
		machines = append(machines, &machine)
	}
	return machines
}

func CreateAdditionMachines() []*agents.AdditionMachine {
	var machines = make([]*agents.AdditionMachine, 0)
	for i := 0; i < config.NUM_MULT_MACHINES; i++ {
		machine := agents.AdditionMachine{
			Id:              i,
			Input:           make(chan agents.MachineWriteOp),
			Logger:          agents.LogChannel,
			IsBroken:        false,
			BreakdownNumber: 0,
			FixMe:           make(chan bool)}
		go machine.RunAdditionMachine()
		machines = append(machines, &machine)
	}
	return machines
}

func RunWorkers() []*agents.Worker {
	workers := make([]*agents.Worker, 0)
	additionMachines := CreateAdditionMachines()
	multiplicationMachines := CreateMultMachines()
	for i := 0; i < config.NUM_WORKERS; i++ {
		outcome := rand.Intn(100)
		var isPatient bool
		if outcome < 50 {
			isPatient = true
		} else {
			isPatient = false
		}
		w := agents.Worker{
			Id:                     i,
			TaskList:               agents.TaskListRead,
			Warehouse:              agents.WarehouseWrite,
			Logger:                 agents.LogChannel,
			MulltMachines:          multiplicationMachines,
			AddMachines:            additionMachines,
			CompletedTasks:         0,
			IsPatient:              isPatient,
			BreakdownReportChannel: agents.ServiceReportWrite,
		}
		go w.Run()
		workers = append(workers, &w)
	}
	return workers
}

func RunCustomers() {
	for i := 0; i < config.NUM_CUSTOMERS; i++ {
		cust := agents.Customer{i, agents.WarehouseRead, agents.LogChannel}
		go cust.Run()
	}
}

func RunService() {
	var ReportCache = make([]agents.BreakdownReport, 0)
	var Reports = make([]agents.BreakdownReport, 0)
	OfficialService := agents.Service{
		Logger:      agents.LogChannel,
		ReportWrite: agents.ServiceReportWrite,
		ReportRead:  agents.ServiceReportRead,
		ReportCache: ReportCache,
		Reports:     Reports,
		FixChannel:  agents.ServiceFixWrite,
	}
	go OfficialService.Run()
}

func RunServiceWorkers(mms []*agents.MultiplicationMachine, ams []*agents.AdditionMachine) []*agents.ServiceWorker {
	workers := make([]*agents.ServiceWorker, 0)
	for i := 0; i < config.NUM_SERVICE_WORKERS; i++ {
		sw := agents.ServiceWorker{
			Id:            i,
			AddMachines:   ams,
			MulltMachines: mms,
			Logger:        agents.LogChannel,
			FixChannel:    agents.ServiceFixWrite,
			ReportChannel: agents.ServiceReportRead,
		}
		workers = append(workers, &sw)
	}

	for _, worker := range workers {
		go worker.Run()
	}
	return workers
}
