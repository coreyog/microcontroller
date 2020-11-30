# Microcontroller

This repo will be the backbone for a PC/Arduino protocol meant to allow a PC to
control a network of Arduino devices.

## Protocol

It's a simple RPC style protocol. Messages are terminated with new lines (no \r)
and must be less than 50 bytes (not including the required new line).

## Useful Packages

 * [board-discovery](https://github.com/arduino/board-discovery)
 * [board-discovery](https://github.com/coreyog/board-discovery) (personal fork: can choose to check serial, network, or both)