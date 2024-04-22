package main

import (
	"fmt"

	"task_20.2.1/pipeline"
)

func main() {
	// Создаем новый пайплайн
	p := pipeline.NewPipeline()
	// Вызываем метод Run и передаем в него нужные методы
	fmt.Printf("Начало работы\nВведите целые числа через пробел или команду stop для остановки конвейера\n")
	p.Run(p.Input, p.StageOne, p.StageTwo, p.StageThree, p.Output)
	// Останавливаем основную горутину
	p.Wait()
	fmt.Println("Работа завершена")
}
