# amore-examples

Examples meant to demonstrate the [Amore Game Framework](https://github.com/tanema/amore)

### test-all

Is used for testing every element of the framework. You can try it out by 
runing `go run main.go` in the directory. It wont look like much but it's 
mostly all there.

### pong 

Pong clone, you can try it out by runing `go run main.go` in the directory.
Use the up and down arrows to play. Enter to restart.

### asteroids

Asteroids clone, you can try it out by runing `go run main.go` in the directory . 
Operate with the arrow keys and space to fire. Destroy the asteroids. Asteroids
also has a pretty good demonstration of the kind of physics you can implement 
yourself.

### physics

Physics demonstrates the usage of [github.com/neguse/go-box2d-lite](https://github.com/neguse/go-box2d-lite)
in amore, for more in depth physics. Normally for most games you dont need this 
level of physics but it's nice to know that it's there.


### racer

A pseudo 3d car racer implemented from the [codeincomplete.com](http://codeincomplete.com/posts/2012/6/22/javascript_racer/)
tutorial. Assets are also from that tutorial.

### platformer

This is a re-implementation/port of [bump.lua](https://github.com/kikito/bump.lua) along with its example program
