package client

type Subject string

const (
	ItemMutateSubject       Subject = "item.mutate.*"
	ItemMutateAddSubject    Subject = "item.mutate.add"
	ItemMutateDeleteSubject Subject = "item.mutate.delete"
	ItemGetSubject          Subject = "item.get.*"
	ItemGetOneSubject       Subject = "item.get.one"
	ItemGetListSubject      Subject = "item.get.list"
)
