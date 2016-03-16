build/check-ecs-agent:
	go build -o build/check-ecs-agent check-ecs-agent.go

clean:
	rm -f build/*
