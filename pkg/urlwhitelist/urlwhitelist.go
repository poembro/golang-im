package urlwhitelist

var Connect = map[string]int{
	// "/pb.ConnectInt/DeliverMessage": 0,
}

// 白名单方法,不需要鉴权的grpc方法写在这里
var Logic = map[string]int{
	// "/pb.LogicInt/SendMessage": 1,
	"/pb.LogicInt/ConnSignIn": 1,
	"/pb.LogicInt/ServerStop": 1,
	"/pb.LogicInt/MessageACK": 1,
	"/pb.LogicInt/Sync":       1,
	"/pb.LogicInt/Offline":    1,
}
