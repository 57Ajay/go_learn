package generics

import "fmt"

func index[T comparable](s []T, x T) int {
	for i, v := range s {
		if v == x {
			return i
		}
	}
	return -1
}

type list[T any] struct {
	next *list[T]
	data T
}

func (l *list[T]) prepend(data T) *list[T] {
	newNode := &list[T]{data: data, next: l}
	return newNode
}

func (l *list[T]) append(data T) *list[T] {
	newNode := &list[T]{data: data, next: nil}
	if l == nil {
		return newNode
	}
	current := l
	for current.next != nil { // Traverse to the end of the list
		current = current.next
	}
	current.next = newNode
	return l
}

func (l *list[T]) printList() {
	current := l
	for current != nil {
		fmt.Printf("%v -> ", current.data)
		current = current.next
	}
	fmt.Println("nil")
}

func (l *list[T]) find(target T) bool {
	current := l
	for current != nil {
		if fmt.Sprintf("%v", current.data) == fmt.Sprintf("%v", target) {
			return true
		}
		current = current.next
	}
	return false
}

func (l *list[T]) isEmpty() bool {
	return l == nil
}

func (l *list[T]) len() int {
	count := 0
	current := l
	for current != nil {
		count++
		current = current.next
	}
	return count
}

func GenericsMain() {
	fmt.Println(index([]int{1, 2, 3}, 2))
	fmt.Println(index([]string{"a", "b", "c"}, "1"))

	var myList *list[int]

	myList = myList.prepend(3)
	myList = myList.prepend(2)
	myList = myList.prepend(1)

	fmt.Println("List after prepends:")
	myList.printList() // Output: 1 -> 2 -> 3 -> nil

	myList = myList.append(4)
	myList = myList.append(5)

	fmt.Println("List after appends:")
	myList.printList() // Output: 1 -> 2 -> 3 -> 4 -> 5 -> nil

	fmt.Println("Length of list:", myList.len())     // Output: 5
	fmt.Println("Is list empty?", myList.isEmpty())  // Output: false
	fmt.Println("Find 3 in list?", myList.find(3))   // Output: true
	fmt.Println("Find 10 in list?", myList.find(10)) // Output: false

	var stringList *list[string]
	stringList = stringList.append("apple")
	stringList = stringList.append("banana")
	fmt.Println("\nString List:")
	stringList.printList()                                                  // Output: apple -> banana -> nil
	fmt.Println("Find 'banana' in string list?", stringList.find("banana")) // Output: true

}
