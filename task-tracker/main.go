package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func saveTask(tasks []Task) error {
	file, err := os.Create("todo.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(tasks)
}

func loadTask(tasks *[]Task) error {
	file, err := os.Open("todo.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(tasks)
}

func AddTask(description string) error {
	var tasks []Task

	if err := loadTask(&tasks); err != nil {
		return err
	}

	newID := 1
	if len(tasks) > 0 {
		newID = tasks[len(tasks)-1].ID + 1
	}

	newTask := Task{
		ID:          newID,
		Description: description,
		Status:      false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tasks = append(tasks, newTask)

	return saveTask(tasks)
}

func ListTasks() error {
	var tasks []Task

	if err := loadTask(&tasks); err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("Nenhuma tarefa encontrada.")
	}

	fmt.Println("Lista de tarefas:")
	for _, task := range tasks {
		status := "[ ]"
		if task.Status {
			status = "[x]"
		}
		fmt.Printf("\tID: %d \n\t\t- Descrição: %s \n\t\t- Status: %s\n", task.ID, task.Description, status)
	}

	return nil
}

func UpdateTask(id int) error {
	var tasks []Task

	if err := loadTask(&tasks); err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = true
			tasks[i].UpdatedAt = time.Now()
			return saveTask(tasks)
		}
	}

	return fmt.Errorf("Tarefa com ID %d não encontrada", id)
}

func DeleteTask(id int) error {
	var tasks []Task

	if err := loadTask(&tasks); err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return saveTask(tasks)
		}
	}

	return fmt.Errorf("Tarefa com ID %d não encontrada", id)
}

func main() {
	cmdAdd := flag.String("add", "", "Adicionar uma nova tarefa")
	cmdList := flag.Bool("list", false, "Listar tarefas")
	cmdUpdate := flag.Int("update", 0, "Atualizar o status de uma tarefa")
	cmdDelete := flag.Int("delete", 0, "Deletar uma tarefa")

	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Uso do Gerenciador de Tarefas:")
		flag.Usage()
		return
	}

	if *cmdAdd != "" {
		if err := AddTask(*cmdAdd); err != nil {
			fmt.Println("Erro ao adicionar a tarefa")
			return
		}
		fmt.Println("Tarefa adicionada com sucesso")
	}

	if *cmdList {
		if err := ListTasks(); err != nil {
			fmt.Println("Erro ao listar tarefas")
			return
		}
	}

	if *cmdUpdate != 0 {
		if err := UpdateTask(*cmdUpdate); err != nil {
			fmt.Printf("Erro ao atualizar a tarefa: %v\n", err)
			return
		} else {
			fmt.Printf("Tarefa %d atualizada!\n", *cmdUpdate)
		}
		return
	}

	if *cmdDelete != 0 {
		if err := DeleteTask(*cmdDelete); err != nil {
			fmt.Printf("Erro ao deletar a tarefa: %v\n", err)
			return
		} else {
			fmt.Printf("Tarefa %d deletada!\n", *cmdDelete)
		}
		return
	}
}
