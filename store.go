package cache

type store interface {
	Add(string,any) error
	Delete(key string) error
	Update(key string,obj any) error
	Get(key string) (any,error)
	List() (any,error)
}
