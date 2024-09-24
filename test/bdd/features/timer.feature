Feature: Timer
  As a gameboy emulator
  In order to sync all the gameboy components
  I want to implement a timer

  Scenario: Static Testing - Registers > check that timer.go file defines a mapping of (register name) -> (register address) with the right key/value pairs
    Then i should have a map "timer_registers" with the following key value pairs:
      | register_name | register_address |
      | DIV           | 0xFF04           |
      | TIMA          | 0xFF05           |
      | TMA           | 0xFF06           |
      | TAC           | 0xFF07           |

  Scenario: Static Testing - Struct > check that a timer instance has the right fields
    Given i instantiate a new struct "Timer" as "timer"
    Then i should have a variable "Timer" with the following fields:
      | field_name | field_type |
      | registers  | map[string]uint8 |
      | divider   | uint8 |
      | counter   | uint8 |
      | modulo    | uint8 |
      | control   | uint8 |
      | cycles    | uint8 |
      | last_cycle | uint64 |

  Scenario: Dynamic Testing - DIV register > check the add subscriber method
    Given i instantiate a new struct "Timer" as "timer"
    Then the timer should have 0 subscribers
    And i add a subscriber
    Then the timer should have 1 subscribers
    And i add 10 subscribers
    Then the timer should have 11 subscribers
