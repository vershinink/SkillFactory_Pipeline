package main

import (
	"log"

	"task_20.2.1/pipeline"
)

func main() {
	// Создаем новый пайплайн
	p := pipeline.NewPipeline()
	// Вызываем метод Run и передаем в него нужные методы
	p.Run(p.Input, p.StageOne, p.StageTwo, p.StageThree, p.Output)
	// Останавливаем основную горутину
	p.Wait()
	log.Printf("Работа завершена")
}
