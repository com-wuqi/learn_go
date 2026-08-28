package main

import "LearnGo/review"

func main() {
	review.RunProtoDemo()
	review.RunUnaryDemo()
	review.RunServerStreamDemo()
	review.RunClientStreamDemo()
	review.RunBidiDemo()
	review.RunMetadataDemo()
	review.RunTimeoutDemo()
	review.RunInterceptorDemo()
	review.RunStatusDemo()
	review.RunRetryDemo()
	review.RunResolverDemo()
	review.RunTlsDemo()
	review.RunKeepaliveDemo()
	review.RunConnStateDemo()
	review.RunStreamInterceptorDemo()
	review.RunSizesDemo()
	review.RunHealthDemo()
	review.RunReflectionDemo()
}
