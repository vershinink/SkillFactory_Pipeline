// Реализация конвейера обработки целых чисел по заданию на основе
// структуры Pipeline. Методы, работающие с потоком целых чисел,
// имеют одинаковые сигнатуры, поэтому могут быть переданы в метод Run
// в любом порядке и количестве.
package pipeline

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"task_20.2.1/ring"
)

// Константы с интервалом и размером буфера
const bufferInterval time.Duration = 2 * time.Second
const bufferSize int = 10

// Pipeline - структурв конвейера, содержит только канал-флаг для остановки
type Pipeline struct {
	stop chan bool
}

// NewPipeline - конструктор конвейера
func NewPipeline() *Pipeline {
	log.Println("Создан новый экземпляр конвейера")
	return &Pipeline{stop: make(chan bool)}
}

// Input - принимает целые числа из консоли, разделенные пробелом.
// Принимет также ключевое слово stop для остановки конвейера. Любые другие
// данные игнорирует. Принимает канал data в аргументах только для
// соответствия сигнатуры остальным методам.
func (p *Pipeline) Input(data <-chan int) <-chan int {
	log.Println("Введите целые числа через пробел или команду stop для остановки конвейера")
	scanner := bufio.NewScanner(os.Stdin)
	output := make(chan int)
	go func() {
		for {
			scanner.Scan()
			input := strings.Fields(scanner.Text())
			for _, v := range input {
				if v == "stop" {
					p.stop <- true
					close(output)
					return
				}
				n, err := strconv.Atoi(v)
				if err != nil {
					continue
				}
				output <- n
			}
		}
	}()
	return output
}

// StageOne - фильтрует отрицательные числа
func (p *Pipeline) StageOne(data <-chan int) <-chan int {
	filtered := make(chan int)
	go func() {
		for {
			select {
			case d, ok := <-data:
				if ok && d >= 0 {
					log.Printf("Стадия 1. Отфильтровано число %d\n", d)
					filtered <- d
				}
			case <-p.stop:
				close(filtered)
				return
			}
		}
	}()
	return filtered
}

// StageTwo - фильтрует числа, не кратные 3, а также 0
func (p *Pipeline) StageTwo(data <-chan int) <-chan int {
	filtered := make(chan int)
	go func() {
		for {
			select {
			case d, ok := <-data:
				if ok && d != 0 && d%3 == 0 {
					log.Printf("Стадия 2. Отфильтровано число %d\n", d)
					filtered <- d
				}
			case <-p.stop:
				close(filtered)
				return
			}
		}
	}()
	return filtered
}

// StageThree - принимает числа и отдает их с определенным в константе интервалом.
// Для этого создает новый кольцевой буфер из локального пакета ring, помещает
// и извлекает из него числа.
func (p *Pipeline) StageThree(data <-chan int) <-chan int {
	filtered := make(chan int)
	r := ring.NewRing(bufferSize)
	go func() {
		var err error
		for {
			select {
			case d := <-data:
				for r.IsFull() {
					continue
				}
				log.Printf("Стадия 3. Отфильтровано число %d\n", d)
				err = r.Write(d)
				if err != nil {
					fmt.Printf("Ошибка работы буфера: %v\n", err)
				}
			case <-p.stop:
				return
			}
		}
	}()
	go func() {
		var d int
		var err error
		for {
			select {
			case <-time.After(bufferInterval):
				for r.IsEmpty() {
					continue
				}
				d, err = r.Read()
				if err != nil {
					fmt.Printf("Ошибка работы буфера: %v\n", err)
					continue
				}
				filtered <- d
			case <-p.stop:
				close(filtered)
				return
			}
		}
	}()
	return filtered
}

// Output - выводит принимаемые числа в консоль. Создает возвращаемый канал
// для соответствия сигнатурам других методов
func (p *Pipeline) Output(data <-chan int) <-chan int {
	output := make(chan int)
	go func() {
		for {
			select {
			case d := <-data:
				log.Printf("Обработанное число: %d\n", d)
			case <-p.stop:
				close(output)
				return
			}
		}
	}()
	return output
}

// Run - метод, запускающий сколько угодно одинаковых по сигнатуре методов конвейера
func (p *Pipeline) Run(stages ...func(<-chan int) <-chan int) {
	log.Println("Запуск конвейера")
	data := make(<-chan int)
	for i := 0; i < len(stages); i++ {
		data = stages[i](data)
	}
}

// Wait - блокирует завершение основной горутины до вызова остановки конвейера
func (p *Pipeline) Wait() {
	<-p.stop
	log.Printf("Завершение работы конвейера... ")
}
