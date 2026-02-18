---
layout: default
title: Home
nav_order: 1
---

# MLGE — Mechanical Lich Game Engine

A simple 2D game engine library built on top of [Ebitengine](https://ebitengine.org/).

## Overview

MLGE is a Go library that provides a collection of game development systems built on top of Ebitengine. Rather than being a standalone engine, it offers reusable packages for common game development needs — entity management, input handling, audio, UI, pathfinding, task scheduling, and more.

## Key Features

- **Entity Component System** — Blueprint-based ECS with component factory and system manager
- **Event System** — Publish/subscribe with immediate and queued dispatch modes
- **Input Management** — Mouse and keyboard input translated into events
- **Audio Playback** — Background music player supporting MP3 and OGG formats
- **A\* Pathfinding** — Allocation-efficient pathfinding with reusable pathfinder instances
- **UI Framework** — Comprehensive widget library with theming and sprite-based rendering
- **Task Scheduling** — Priority-based task system with proximity assignment and escalation
- **State Machine** — Stack-based game state management
- **Sensory Simulation** — Grid-based stimulus propagation for sound, scent, and pheromones
- **Resource Management** — Asset loading and caching for textures, fonts, and sounds
- **Text Rendering** — Text drawing, measurement, and word wrapping
- **Dice Rolling** — Tabletop-style dice expression parser
- **Utility Functions** — Math, geometry, and drawing helpers

## Packages

| Package | Description |
|---------|-------------|
| [`ecs`](ecs.md) | Entity Component System with blueprints and system manager |
| [`event`](event.md) | Publish/subscribe event system |
| [`input`](input.md) | Mouse and keyboard input management |
| [`audio`](audio.md) | Background music playback (MP3/OGG) |
| [`path`](pathfinding.md) | A* pathfinding |
| [`resource`](resource.md) | Asset loading and caching |
| [`state`](state-machine.md) | Stack-based state machine |
| [`task`](task.md) | Priority-based task scheduling |
| [`sense`](sense.md) | Grid-based sensory simulation |
| [`ui/minui`](ui.md) | Immediate UI framework with widgets and theming |
| [`text`](text.md) | Text rendering and wrapping |
| [`dice`](dice.md) | Dice expression parser |
| [`message`](message.md) | In-game message log |
| [`utility`](utility.md) | Math, geometry, and drawing helpers |

## Installation

```bash
go get github.com/mechanical-lich/mlge
```

## Dependencies

MLGE is built on top of [Ebitengine v2](https://ebitengine.org/), a simple 2D game engine for Go. Ebitengine handles the low-level rendering, audio context, and input polling, while MLGE provides higher-level game systems on top of it.

## License

See [LICENSE](https://github.com/mechanical-lich/mlge/blob/main/LICENSE) for details.
