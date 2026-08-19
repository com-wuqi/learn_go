package main

import "LearnGo/review"

func main() {
	review.RunProtoDemo()
	review.RunUnaryDemo()
	review.RunServerStreamDemo()
	review.RunClientStreamDemo()
	review.RunBidiDemo()
}
