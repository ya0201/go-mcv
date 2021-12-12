package nozzle

import 

type Nozzle interface {
  Pump(chan Comment) error
}
