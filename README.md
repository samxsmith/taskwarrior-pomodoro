# Taskwarrior Pomodoro

A simple pomodoro tool for task warrior.

## Requirements
- Go installed, with `$GOPATH` variable set
- Taskwarrior (as `task` bin)

## Install
```
git clone git@github.com:samxsmith/taskwarrior-pomodoro.git
make install
```

## Use
```
$ pomdoro
For how many minutes?: 25
Which taskwarrior task ID?: 1
You want to run task: My Taskwarrior Task

y/N: y
```

The timer will run for the number of minutes you specify.
It will mark your taskwarrior task as started.

At the end of the interval, it will sound an alarm, and stop the task.

To exit, `Ctrl-C` and the task will be stopped.
