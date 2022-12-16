package exception

// 存储声明的ErrorCode
var store = &Stroe{}

type NamespaceStroe struct {
	Namespace    string
	Excetiptions []APIException
}

type Stroe struct {
	items []NamespaceStroe
}

func (s *Stroe) Items() []NamespaceStroe {
	return s.items
}

func (s *Stroe) Add(e APIException) {

}
