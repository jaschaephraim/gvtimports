package main

type set struct {
	els map[string]interface{}
}

func newSet() *set {
	return &set{
		els: map[string]interface{}{},
	}
}

func (s *set) add(element string) {
	s.els[element] = nil
}

func (s *set) elements() []string {
	els := make([]string, len(s.els))

	i := 0
	for k := range s.els {
		els[i] = k
		i++
	}
	return els
}
