package envy

import (
	"log"
	"reflect"
	"runtime"
	"sync"
)

type Object interface{}
type Slot interface{}
type SlotVal reflect.Value
type Signal struct {
	id int
}

var idCounter int = 0
var idMutex sync.Mutex = sync.Mutex{}
var registeredSlot map[int][]SlotVal = map[int][]SlotVal{}

/* Create a new signal */
func New() Signal {
	idMutex.Lock()
	signal := Signal{idCounter}
	idCounter++
	idMutex.Unlock()
	return signal
}

/* Connect a signal to a slot */
func Connect(signal Signal, slot Slot) {
	if !isFunction(slot) {
		log.Println("error: not a slot/function")
		return
	}

	signalID := signal.id

	if registeredSlot[signalID] == nil {
		registeredSlot[signalID] = []SlotVal{}
	}

	slotval := SlotVal(reflect.ValueOf(slot))
	if !containsSlot(registeredSlot[signalID], slotval) {
		registeredSlot[signalID] = append(registeredSlot[signalID], slotval)
	} else {
		log.Println("warning: slot already connected to signal")
	}
}

/* Deletes all slots registered to a signal */
func DeleteSlots(signal Signal) {
	delete(registeredSlot, signal.id)
}

/* Emit a signal with the given parameters */
func Emit(signal Signal, params ...Object) {
	signalID := signal.id

	if registeredSlot[signalID] == nil {
		registeredSlot[signalID] = []SlotVal{}
	}

	slots := registeredSlot[signalID]

	for _, slot := range slots {
		slottype := reflect.Value(slot).Type()
		paramlist := []reflect.Value{}
		/* Try to match number of parameters */
		if slottype.IsVariadic() {
			for _, param := range params {
				paramlist = append(paramlist, reflect.ValueOf(param))
			}
		} else {
			for i := 0; i < slottype.NumIn(); i++ {
				paramlist = append(paramlist, reflect.ValueOf(params[i]))
			}
		}
		go func() {
			defer crashHandler("Fatal: Cannot call slot/function")
			reflect.Value(slot).Call(paramlist)
		}()
	}
}

func containsSlot(slots []SlotVal, slot SlotVal) bool {
	for _, s := range slots {
		if slot == s {
			return true
		}
	}
	return false
}

func isFunction(obj Object) bool {
	f := reflect.ValueOf(obj)
	return f.Kind() == reflect.Func
}

func crashHandler(msg string) {
	if r := recover(); r != nil {
		log.Println(msg)
		log.Println(r)
		// Print stack
		bufsize := 4096
		buf := make([]byte, bufsize)
		runtime.Stack(buf, false)
		log.Println(string(buf))
	}
}
