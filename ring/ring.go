// Реализация кольцевого массива на основе среза указателей, так как
// указатели удобно проверять на nil при использовании в методах,
// где нужно вычислить свободное место в буфере
package ring

import (
	"fmt"
	"sync"
)

// Ring - кольцевой массив
type Ring struct {
	data []*int
	m    sync.Mutex
	size int
	head int
	tail int
}

// NewRing - конструктор кольца
func NewRing(size int) *Ring {
	return &Ring{
		data: make([]*int, size),
		m:    sync.Mutex{},
		size: size,
		head: 0,
		tail: 0,
	}
}

// IsFull - проверка, есть ли место в буфере для записи
func (r *Ring) IsFull() bool {
	return r.data[r.tail] != nil
}

// IsEmpty - проверка, есть ли значения в буфере для чтения
func (r *Ring) IsEmpty() bool {
	return r.data[r.head] == nil
}

// Write - записывает значение в буфер
func (r *Ring) Write(n int) error {
	if r.data[r.tail] != nil {
		return fmt.Errorf("буфер полон")
	}
	r.m.Lock()
	defer r.m.Unlock()
	r.data[r.tail] = &n
	r.tail++
	if r.tail == len(r.data) {
		r.tail = 0
	}
	return nil
}

// Read - вынимает значение из буфера
func (r *Ring) Read() (int, error) {
	if r.data[r.head] == nil {
		return -1, fmt.Errorf("буфер пуст")
	}
	r.m.Lock()
	defer r.m.Unlock()
	n := r.data[r.head]
	r.data[r.head] = nil
	r.head++
	if r.head == len(r.data) {
		r.head = 0
	}
	return *n, nil
}
