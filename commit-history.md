* commit f2aeade4aeee427e1a4cc3e9d892c17d3315cbdc
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Feb 17 23:01:15 2025 +0100
| 
|     Debugger: reenabling the /debugger route // Adding Tick & Run funcs // removing all state channels except for cpu & memoryWrites (work in progress)
| 
* commit e2797be0bdd0ee287b342b57035d3f25ae573219
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Feb 16 23:23:19 2025 +0100
| 
|     debugger: Define a Tick() and a Run() funcs to allow piloting the debugger: the debugger make no use of the gameboy state action channel // debugger always defines an internal cpu state channel even if the client/caller does not since cpu.pc is necessary to handle breakpoints
| 
* commit 41eab01ff3802537f7d088c6e086c2bc127f3345
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Feb 16 23:20:56 2025 +0100
| 
|     bus: reenable memory writes whenever a write operation is performed on memories // gameboy: add a public Tick() func for the debugger // send the initial state after loading a rom or resetting the gameboy
| 
* commit 33722f6a192ec81a1dea54b7564d30dd0f59ab12
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Feb 16 18:48:49 2025 +0100
| 
|     go_project_structure.md: explain go mono repo with master module and submodules with multiple packages
| 
* commit 9394eced6c8d3226794c5d2e5c3880b689f5cb1a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Feb 15 23:06:56 2025 +0100
| 
|     execution_flow_chart.md: documenting the data structure package and the Iterable interface
| 
* commit 5c31578d59addba01b0107f63b973c6da50b22da
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Feb 15 23:04:41 2025 +0100
| 
|     Refactoring: all iterations and transformation capabilities are now decoupled from fifo data structure using Iterable[T] interface and can now be reused in further data structures (lifo, stacks, ...)
| 
* commit 707d767a6a5bd927c3bd0df9429608a7cf7d19cf
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Feb 15 20:40:35 2025 +0100
| 
|     fifo: adding a mapper func that applies a func fn(*T) *U to  all fifo values and returns the  new array over a channel of type *U
| 
* commit b596bb72d83806f20edc0898033153c51944db4d
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Feb 14 05:32:00 2025 +0100
| 
|     added a file to document the generation of execution flow charts
| 
* commit 9ccacd8a78d3172caeae1f4809b164fdf5435291
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Feb 14 05:31:31 2025 +0100
| 
|     created a new package debugger in which the debugger is now located // created a new execution_flow data struct to save any kind of state during execution and allow to extract data from it to generate for ex an pc addresses execution flow chart (work in progress)
| 
* commit d6e41d0025da3251c931dfc36b6337c58bbe855d
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Feb 14 05:29:46 2025 +0100
| 
|     created a new package datastructure to hold all helper data struct used in the debugger: node, fifo, iterator
| 
* commit fcf0896cae3d2f3e3dddaabba181c9c94c28caa4
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Feb 14 05:28:34 2025 +0100
| 
|     removed BDD test cases & cucumber dependencies // moved the debugger to its own package
| 
* commit b35af99791e6e9b5ec778a7921973d78b9292fce
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Feb 12 07:47:39 2025 +0100
| 
|     bus: memory writes are disabled for now (used in debugger but since i am not using them they are not reset and end up occupying huge memory space)
| 
* commit ff08231de6388110a6b5a4b24e6f659eb19ceb0b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Feb 12 07:46:19 2025 +0100
| 
|     bus: read to joypad register 0xFF00 should return 0xFF since it joypad is not implemented yet otherwise 0x00 signifies that all buttons are being pressed which resets some games like tetris // reads to prohibited memory zone FEA0-FEFF should return 0xFF (apparently tetris has some bugs among which trying to access this forbiddent zone) // writes to forbidden zone 0xFEA0-FEFF should be ignored (once again a bug in tetris)
| 
* commit 05916ff272acfa8a1f779b0ef833c5574ccaebc1
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Feb 12 07:38:48 2025 +0100
| 
|     cpu interrupts: handle VBLANK interrupt and don't forget to reenable ime flag after handling any interrupt on the next instruction cycle (during cpu state CPU_EXECUTION_STATE_EXECUTE)
| 
* commit 9da275408ac33764a89026a8c96debff299fea35
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Feb 12 07:34:49 2025 +0100
| 
|     ppu: request VBLANK interrupt at the end of each frame (144,0)
| 
* commit 99f64979c177805f7b653ad279a5f3f1f62bc8ce
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Feb 12 07:32:45 2025 +0100
| 
|     gameboy.run(): handle cpu interrupts when cpu is in state 'fetch'
| 
* commit 0d805762f1de9485dbba50fdb7faf139f8078bbc
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Feb 6 11:20:13 2025 +0100
| 
|     ppu: the background is now correctly drawn to screen :)
| 
* commit e4f8539dff1ece9f9646a8a507dd1b0b8dbff03b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Feb 6 11:19:16 2025 +0100
| 
|     memory: added ResetWithZeros() func to reset a memory with all zeros
| 
* commit 77a17cbf9fb7e54f6185914078697d8989f884be
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Feb 6 11:18:26 2025 +0100
| 
|     gameboy: i am not storing the tick duration in a variable before using it in time.Sleep() to gain some time and reach a better frame rate
| 
* commit 1293eef3be9f155a72d0ab27069dd907d891f9aa
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Feb 6 11:16:49 2025 +0100
| 
|     CPU: HRAM, I/O registers and IE will now be initialized and reset with 0 data (otherwise I think I have not reproducible gameboy boot sequence
| 
* commit 0b8a60d897a03fe8ae63962c2f82d5dab8164767
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Feb 3 09:41:10 2025 +0100
| 
|     joypad.go: added the logic to handle button presses along with associated registers constants // removed the joypad channel (TODO: check if this is a good choice ...)
| 
* commit d749776c7905dace221bc33147e0b43cc647b25b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Feb 3 09:39:39 2025 +0100
| 
|     main.go: same as previous commit: commenting all lines that do not work since the heavy changes in gameboy interface and gameloop
| 
* commit 89a7f93ec97e9329c858c34492c812b085046118
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Feb 3 09:38:37 2025 +0100
| 
|     debugger: it needs to be refactored since the gameloop and gameboy interface was heavily changed: for now, just commenting the 2 lines causing issues and not using /debugger endpoint for now
| 
* commit 5ec70106b251c834c9307b6bdf60f73ecc06f708
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Feb 3 09:36:56 2025 +0100
| 
|     ppu.go: fixed the background rendering routine who had miscalculations all over the place: wrong screen to draw x position, wrong tile id & data, wrong pixel data
| 
* commit d5365b0b04334aeb1de956b504bf76fe9dacd2db
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Feb 3 09:35:15 2025 +0100
| 
|     the ppu.getState() interface changed: it now returns the state along with an error
| 
* commit f5c709486a18ad0a0443eb2eec07736afc7e5408
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Feb 3 09:33:12 2025 +0100
| 
|     gameboy.go: the tickDuration inside the run gameloop is now calculated after reporting the gameboy state through the state channels to account for that and not wait too long // the get ppu.getState() interface changed: it nows returns the state along with an error // ACTION_LOAD_GAME now sets the gameboy state to GB_STATE_PAUSED (otherwise the state machine is blocked in state GB_STATE_NO_GAME_LOADED
| 
* commit a209ba240ef8aba20096d7102f995c22043ff36d
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Feb 1 23:17:57 2025 +0100
| 
|     debugger.go: remove public funcs Tick, Run, Pause, Resume and Stop which are now managed through gameboyActionChannel by the gameboy directly // adding gameboyActionChannel <-chan GameboyActionMessage // note: this component will be checked when testing the /debugger endpoint from the server after making the default /gameboy endpoint work
| 
* commit 2c916b64581559ff6fd17c52e3e067ee4996c363
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Feb 1 23:15:18 2025 +0100
| 
|     gameboy.go: refactor stateMachineListener to use for/select channel construct to listen to the gameboyActionChannel // remove the public Tick func which is not used anymore (should add a GB_ACTION_TICK to handle this case and call tick() instead // refactor sendState() func to only report state if the corresponding channel is defined (cpu,ppu,apu,memorywrites)
| 
* commit 217ae22eece0780c118ea784f6ad4acc1a742f5f
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Feb 1 21:42:37 2025 +0100
| 
|     updating the way the gameboy runs: it now used a state machine to manage its execution through a State Action channel than operates the gameboy state MDP model and drives the run function
| 
* commit af483e9c30fdda78d2f84303a3a3fa5b4acdb7cb
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Feb 1 21:40:12 2025 +0100
| 
|     new document about the gameboy state: explains how the gameboy state machine works // updated the gameboy architecture document
| 
* commit 27759ce990950a455f7cb1ffbd0097f1c89a294f
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Jan 30 23:20:50 2025 +0100
| 
|     adding a file about the architecture: will contain all explanations about the link between structs along with diagrams and explanation of the main features of each element of the gameboy without going into the details, just the struct and messages exchanges
| 
* commit 446d05c2e83d02cdf8e907ed94a32d803fc4241f
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Jan 10 06:27:26 2025 +0100
| 
|     lib_test: remove mmu dependencies
| 
* commit 049d9a0865f316ec8fb035765ec49ad6d77ed733
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Jan 10 06:26:29 2025 +0100
| 
|     update docs
| 
* commit 751c71ead831a909606ed6b12bfd7c9b4154b0f3
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Jan 10 06:26:03 2025 +0100
| 
|     gameboy & debugger: remove joypad channel // access memory directly through bus and remove mmu dependencies
| 
* commit dfa355bee9a038ff34b3549f794218dbed0ec76c
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Jan 10 06:24:08 2025 +0100
| 
|     ppu: remove cpu dependency and use bus directly to access memory
| 
* commit 2fb380aad82427fde05ac96a20aa1b7e9d20d426
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Jan 10 06:22:16 2025 +0100
| 
|     merge mmu.go & bus.go into bus.go: bus was only calling mmu but not adding any func to it
| 
* commit 2519c497f40a013b5d0e3e00e65758b66f096524
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Jan 10 06:20:50 2025 +0100
| 
|     test suites: rename files and delete unused ones
| 
* commit 01c954a30174de8e07da7e8ce3dcf4aa8ba8123c
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Jan 3 12:21:03 2025 +0100
| 
|     APU: adding globals to reference special APU registers // renaming the onTick func to Tick
| 
* commit 8a75171eb81166ef6e7cabd54f5cc3d4cd5e59e4
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Jan 3 09:57:04 2025 +0100
| 
|     adapting the gameboy & debugger to use the tick func to advance the whole gameboy state (cpu, ppu, apu) by one clock cycle and return the state to the front-end after a tick or when an image is ready on run
| 
* commit add16f1b80abacfd365b62a71ddf662b77c7b052
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Jan 3 09:55:22 2025 +0100
| 
|     PPU: adding implementation for the rendering of the background as a state variable image (256x256x2bp)
| 
* commit ce22bfb4bb4a31ccd19c939fa681d95bb6b80596
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Jan 3 09:53:40 2025 +0100
| 
|     RET instruction wasn't updating the cpu cycles after execution
| 
* commit 418324d4319ae1fae6cc01d65c623206090ad8f3
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Dec 31 13:17:59 2024 +0100
| 
|     cpu registers now uses the new naming convention for IE register
| 
* commit 7bf553611a9aa6646335e11863b5db89aa514fee
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Dec 31 13:04:35 2024 +0100
| 
|     debugger & gameboy test suites containing test using the old cpu Step & Run func are disabled now and should be refactored & maybe reworked
| 
* commit 209745085e466135ffdd8887214144c237ef58c6
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Dec 31 13:03:03 2024 +0100
| 
|     adapt the cpu test suites to the new changes made to the cpu logic
| 
* commit 27d059f76c05ff644f8be1d2149964e093f6924e
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Dec 31 13:01:16 2024 +0100
| 
|     MAJOR CHANGE: to get towards a time accurate implementation and have the cpu & ppu work sync more realistically, the cpu is not anymore executing a full instruction in one tick but instead cycles through fetch, decode, execute + stall (to let the gameboy clock align with instructions processing time). In consequence, Step and Run are now replaced by a tick function that selects the right action (CPU_EXECUTION_STATE)
| 
* commit 0c5d7931f41c89e9b6b9c775b4c9bb2458dba239
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Dec 31 12:49:31 2024 +0100
| 
|     timer: remove the Synchronizable interface defining the onTick() func since the Timer doesn't use anymore the Subscriber design pattern to drive the gameboy components (cpu, apu, ppu) but the components are instead driven by the gameboy itself which checks the processing time between 2 ticks to determine how much time to stall before triggering the next tick inside a classical 'game loop'
| 
* commit abe9d8735fb96a80c3e401a68844fd136b21d194
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Dec 26 13:04:41 2024 +0100
| 
|     Timer: adding implementation and test suite for the gameboy timer (manage DIV, TIMA, TMA, TAC registers and TIMA overflow interrupt logic // MMU: when writing to div register with Write func @0xFF04, reset the DIV register value // MMU: adding a timerWrite func that allows the Timer to write to the DIV reg any value to allow it to record the current tick @16,384Hz
| 
* commit 49c19beff7265e8e771d678608d93134cf300405
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Dec 18 09:04:39 2024 +0100
| 
|     STOP instruction should reset DIV register @0xFF04: added implementation and updated corresponding TC
| 
* commit 4da97a6a36bfc61b8f006e13a19732f562993a8d
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Dec 16 09:47:50 2024 +0100
| 
|     refactoring: moved cpu.onTick from cpu_interrupts.go to cpu.go to isolate interrupt handlers & because it makes more sense as onTick is related to cpu.Step func which is already in cpu.go
| 
* commit 9f4e302daf7e60e3787d9d8e9fd3290dd921b5a2
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Dec 16 09:42:07 2024 +0100
| 
|     adding interrupt handlers for VBLANK, LCD STAT, Timer, Serial & Joypad. Basically resets ime & if registers, push PC to SP, waits for 5 M-cycles and finally jumps to interrupt vector , , ,  or
| 
* commit eba77405f7462288f4f1fdea05423a1a6364f063
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Dec 10 09:47:57 2024 +0100
| 
|     unmap bootrom from mmu on write to 0xFF50 register
| 
* commit cd4ac211cee7ee5ffa03a94d407394679f42e160
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Dec 10 00:19:14 2024 +0100
| 
|     separate the logic to initialize a gameboy (link all the components together) and load a new rom (reset all the components to a random state like when the gameboy is turned on) allowing to load a rom, then another one
| 
* commit 6f6f7430956b5df75b93dc4281733ef34d9442c9
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Dec 9 23:44:28 2024 +0100
| 
|     DEC instruction: forgot else statements when checking H & Z flags for H & [HL] operands causing the test to randomly fail
| 
* commit 66ee64ddffdd07c0eaed6ee31bae1804b4aa3181
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Dec 9 22:10:39 2024 +0100
| 
|     moved all gameboy initialisation code from gameboy_init.go to gameboy.go (no real reason to split the file since it is part of the LoadRom implementation) // moved the bootrom initialization from cpu.go to gameboy.go because it makes more sense to have it as a feature bound to the gameboy architecture itself and not the cpu even if it is physically located in a rom inside the cpu SOC ... anyway i can still change it if necessary
| 
* commit 42c6d08bc117c56faf17c11a103c3daa83cd9172
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Dec 9 22:02:58 2024 +0100
| 
|     gameboy_state.go only contained support funcs that printed out the cpu state to the console but were not used --> file deleted
| 
* commit 505fe394e1e914bb1c4139847eab5d805d2f6fcf
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Dec 9 22:01:21 2024 +0100
| 
|     renaming test files with 0_ prefix to separate them from the source files in the gameboy folder for clarity
| 
* commit 787bdab1000de2c8238db480f3307c9ef274373b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Dec 8 12:30:46 2024 +0100
| 
|     SCF: simplify logic by using cpu.setCFlag() func instead of logic OR operation on register F
| 
* commit 3f353d53e2dec2c14be688b5be8f6ac75aa2a03b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Dec 8 12:00:01 2024 +0100
| 
|     DEC instruction: flags H and Z were not reset properly (missing else statement in instruction handler for all r8 opcodes // adapted the DEC test suite to randomize flags before executing the test program
| 
* commit 7557239946bc67e5d9a52512a8c462e60f776a2e
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Dec 7 13:24:19 2024 +0100
| 
|     transmit concrete values through channels instead of pointers to be sure that by the time they are transmitted trough the server to the front-end, these values still exist and are not update leading for example a *[]MemoryWrite to point to a struct with a different size (5 updates in the previous step and only 2 for ex, leading to issue when marshalling because we get an out of range issue)
| 
* commit e674152046411ef44af15be3304a503a0dd8bf28
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Dec 7 12:57:26 2024 +0100
| 
|     clearing memory writes in Gameboy.onTick() func to avoid concurrent access causing a bug leading the gameboy to panic while waiting for VBLANK @bootrom 0x60-0x68
| 
* commit 65ba422084fbe97898fb5080e9ba4eb44f068569
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Dec 6 23:22:11 2024 +0100
| 
|     cpu state's instruction cycles is of type []uint8 and is not transmitted correctly because it should used a custom marshaller when encoded by the gameboy-go server: therefore, the cycles attribute of the Instruction struct should be changed to JSONableSlice
| 
* commit 25ca4b5cf980a4bad0ef9c6fd5e8fa971e0b80f6
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Dec 5 21:11:29 2024 +0100
| 
|     adding implementation for instruction CB SET 0xC0 - 0xFF // adding test suite for instruction CB SET
| 
* commit 889043e6433625411f12db559dcdd5a6d2bf0b2b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Dec 5 20:25:17 2024 +0100
| 
|     adding implementation for instruction CB RES 0x80 - 0xBF // adding test suite for instruction CB RES
| 
* commit ea1ceadd2543e47622fdbc68e5c2dca030576acf
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Dec 5 19:14:13 2024 +0100
| 
|     adding implementation for instruction CB BIT 0x40 - 0x7F // adding test suite for instruction CB BIT
| 
* commit b65096af54a6848a02c81595fdc8a8d7030cb97c
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Dec 5 18:25:48 2024 +0100
| 
|     adding implementation for instruction CB SRL 0x38 0x39 0x3A 0x3B 0x3C 0x3D 0x3E 0x3F // adding test suite for instruction CB SRL
| 
* commit c2e14e0b14196be2af1ee8bf0153611ab1a132b9
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Dec 5 12:45:03 2024 +0100
| 
|     adding implementation for instruction CB SWAP 0x30 0x31 0x32 0x33 0x34 0x35 0x36 0x37 // adding test suite for instruction CB SWAP
| 
* commit 05b529cee7461fd33d54beb2febd5c2806c19543
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Dec 5 10:11:37 2024 +0100
| 
|     adding implementation for instruction CB SRA 0x28 0x29 0x2A 0x2B 0x2C 0x2D 0x2E 0x2F // adding test suite for instruction CB SRA
| 
* commit 507da0e5e8f47e046d30ffafd11c4dc6539249d7
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Dec 5 08:12:31 2024 +0100
| 
|     adding implementation for instruction CB SLA 0x20 0x21 0x22 0x23 0x24 0x25 0x26 0x27 // adding test suite for instruction CB SLA
| 
* commit 1af9402c6e4cc7feef0ae7afbb608bd9f5692e7b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Dec 2 22:52:57 2024 +0100
| 
|     adding implementation for instruction CB RR 0x18 0x19 0x1A 0x1B 0x1C 0x1D 0x1E 0x1F // adding test suite for instruction CB RR
| 
* commit cb92cdb5c1ceba81fbfc5bcc3c67d1753f4f382a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Dec 2 22:37:18 2024 +0100
| 
|     adding implementation for instruction CB RL 0x10 0x11 0x12 0x13 0x14 0x15 0x16 0x17 // adding test suite for instruction CB RL
| 
* commit 0e4de0cc03ab8c4328ab1fe74647fec23cdf879e
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Dec 2 22:13:46 2024 +0100
| 
|     adding implementation for instruction CB RRC 0x08 0x09 0x0A 0x0B 0x0C 0x0D 0x0E 0x0F // adding test suite for instruction CB RRC
| 
* commit ada2bff7d1fc6c5fc97d12b8f8bf946583d8050c
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Dec 2 21:05:13 2024 +0100
| 
|     adding implementation for instruction CB RLC 0x00 0x01 0x02 0x03 0x04 0x05 0x06 0x07 // adding test suite for instruction CB RLC
| 
* commit beac79d3cd4e71a2418bc42a1ee0ade94e89cc8c
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Dec 2 21:04:19 2024 +0100
| 
|     moving helper func and global vars used by cpu testing files (cb & non-cb) to lib_test.go
| 
* commit cb38843c63c357a295d9490ef890634129d3c4af
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Dec 2 08:59:29 2024 +0100
| 
|     adding implementation for last non CB instruction AND 0xA0 0xA1 0xA2 0xA3 0xA4 0xA5 0xA6 0xA7 0xE6 // adding test suite for AND instruction
| 
* commit 59318f0d90101ccd4a62848083de3fc873416298
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Dec 2 08:31:22 2024 +0100
| 
|     adding implementation for instruction ADC 0x88 0x89 0x8A 0x8B 0x8C 0x8D 0x8E 0x8F 0xCE // adding test suite for ADC instruction
| 
* commit dbeae78aa0c55b5f3ae8b34e6d17244f2d6a8f85
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Dec 1 10:44:19 2024 +0100
| 
|     adding implementation for instruction ADD 0x80 0x81 0x82 0x83 0x84 0x85 0x86 0x87 0xC6 0xE8 // adding test suite for ADD instruction
| 
* commit 1e6c97135b50571f33aa667763ee6e8baf000c3e
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Nov 30 13:05:28 2024 +0100
| 
|     adding test suite for instruction SBC
| 
* commit 608e4a01633270d00fb1253a41b33dda4263aecc
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Nov 30 13:05:12 2024 +0100
| 
|     adding implementation dor instruction SBC 0x98 0x99 0x9A 0x9B 0x9C 0x9D 0x9E 0x9F 0xDE
| 
* commit a359b0c72e1813a9a84941b8f2fe786d791d414c
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Nov 28 14:07:01 2024 +0100
| 
|     adding test suite for instruction SUB opcodes 0x90 0x91 0x92 0x93 0x94 0x95 0x96 0x97 0xD6
| 
* commit 311bc2989bf8f857e165bf6dd5924a38cdf3a48d
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Nov 28 14:05:51 2024 +0100
| 
|     instruction SUB refactoring (use cpu.operand instead of switch case on operand) // fixed a bug that made opcode 0x96 SUB A, [HL] not supported (naturally solved by using cpu.operand that by design correctly fetches the indirect [HL] operand value)
| 
* commit 098945d38f606f3d8af0deb00094ab670d1e504f
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Nov 28 10:12:57 2024 +0100
| 
|     added test suite for CPL instruction
| 
* commit 304ab2462db5f35aaff53f19f6388691385d076d
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Nov 28 00:07:42 2024 +0100
| 
|     adding implementation for instruction OR 0xB0 0xB1 0xB2 0xB3 0xB4 0xB5 0xB6 0xB7 0xF6 // adding test suite for instruction OR
| 
* commit cd2105a3892286985d05ce40b7dbe7f1f9074c74
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Nov 27 21:04:12 2024 +0100
| 
|     adding test suite for POP instruction
| 
* commit 914020263bdc8f822011ba2d8a2598188dd4873f
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Nov 27 21:03:55 2024 +0100
| 
|     refactoring POP instruction implementation by using func cpu.pop()
| 
* commit 7dfcf1c2155bdbbc9aba7da7255363ba9b64777b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Nov 27 20:16:30 2024 +0100
| 
|     removing one more test file that should be refactored
| 
* commit 25a1873da52b82eb4a35c29047ed78f6b5e7e38c
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Nov 27 19:30:56 2024 +0100
| 
|     stop tracking temporary test files
| 
* commit 4aa14d6f66ff86c3f218f6273867c83c0e3ceeb6
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Nov 27 19:15:13 2024 +0100
| 
|     disabling 2 buggy tests to prevent github test action to always fails, thus not providing any valuable information (these tests were not that intersting anyway and needed to be refactored)
| 
* commit 59fbed9fd39c13b520a35a5830a18fd2ec692575
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Nov 27 19:13:35 2024 +0100
| 
|     adding test suite for instruction PUSH
| 
* commit 511bf6b254b1ef249d99af838fc254d1e1582c32
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Nov 27 19:13:02 2024 +0100
| 
|     simplifying instruction PUSH implementation by using cpu.push func that already does the job, thus avoiding to have to maintain 2 different versions
| 
* commit e9c33fcca9fa779fe756381eada38f53cfeee249
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Nov 27 19:11:35 2024 +0100
| 
|     adding support for operand AF immediate used in PUSH instructions 0xC5 0xD5 0xE5 0xF5
| 
* commit 7a8827111b53470210767f62aa74229d9a51ba75
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Nov 27 13:01:44 2024 +0100
| 
|     adding test suite for LD instruction // corrected an issue with HL operand not being in/decremented when fetched by the cpu in opcodes 0x22 0x2A 0x32 0x3A involving [HL+/-] operand // fixed opcode 0xF8 LD HL, SP + e8 and flag calculation
| 
* commit 113d22f7d10fb4576867e3e8ebba92e599061da3
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Nov 23 12:00:43 2024 +0100
| 
|     adding test suite for instruction CP
| 
* commit c7e252ac87f12c5b9c10f3316fe36edfecd9a946
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Nov 23 09:29:53 2024 +0100
| 
|     added test suite for DEC instruction // corrected a bug for 0x35 DEC [HL] where flags were not set at all // added a func to randomize flags in test cases
| 
* commit 78b2474f4ff55240202556cbc97273a27275a32b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Nov 22 23:34:31 2024 +0100
| 
|     adding test suite for xor instruction
| 
* commit 86b25ac4ad1d86cf95754a62f82b2c187dc9cf19
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Nov 20 13:19:59 2024 +0100
| 
|     INC instruction: adding support for operands HL immediate, BC, DE & SP // adding corresponding test cases in test suite TestINC
| 
* commit 83ee5152e62fa255cfbcc4a8d2345c49d845b80e
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Nov 18 15:38:29 2024 +0100
| 
|     implementing and using a fifo data structure for message channels (program flow/cpu/ppu/apu/memory/joypad states) to make sure the updates come in order and to make the channels non-blocking
| 
* commit 145460b0209bd86be0f61afcd07385c56fba1f6a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Nov 18 15:35:16 2024 +0100
| 
|     refactoring long tests (JR/JP) that were too long into a test suite with a separate test func for each opcode // adding a test suite for CALL instruction // harmonizing the error output by naming the test in error within a test suite
| 
* commit 46f354ca4bae0c3336451fa8127898533a579b56
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Nov 18 15:33:11 2024 +0100
| 
|     updating the opcodes.json file to differentiate between the 'C' register and flag register // updating the cpu fetch operand func accordingly // update the cpu instructions handlers accordingly to use flag_Z & flag_C instead of C & Z to make it more obvious which operand we are using in the instruction
| 
* commit 9d0b798cbeaa81929fc32a0005e5e0491955a50b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Nov 18 15:26:48 2024 +0100
| 
|     deleting irrelevant documentation used during the dev process
| 
* commit 33a8b8247f5a059ea46748cd77eabb7a41ce8afb
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Nov 18 15:25:34 2024 +0100
| 
|     added a document explaining how we can regroup the implementation of several opcodes based on their mnemonic and how we should fetch the operands in a systematic way // it should also highlight the changes made to opcodes.json file to distinguish the 'C' register from the 'C' flag
| 
* commit c7fc0cffb09f9f340d4b55d5efa0914cffe3eb1a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Nov 18 15:21:57 2024 +0100
| 
|     adding an image of the ui to the blog md file to explain how i debugged the bootrom sequence execution
| 
* commit e96b021fbc03af13f5a0656c391ec5d5aea5505f
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Nov 18 15:18:43 2024 +0100
| 
|     added updating gameboy go core library to explain how to make the changes in gameboy go core module available to the gameboy server
| 
* commit f66754b8d4d65ff17b28526af5bb9e5bf7a11c05
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Nov 18 15:17:34 2024 +0100
| 
|     added ppu.md to list explain how the ppu is working (layers & registers)
| 
* commit 9518ab58996193c5e5c3a638034ba515819c1877
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Nov 18 15:15:57 2024 +0100
| 
|     moved blog.md tp docs folder
| 
* commit 89c542482112e3b5a154ca24f1ae8676bc8a67f4
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Oct 9 01:40:46 2024 +0200
| 
|     CP: flags H & N were not set correctly
| 
* commit b81748ad9fd0286d4125301607275b8dd81ea3f3
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Oct 9 00:28:17 2024 +0200
| 
|     sub instruction: i forgot to update the pc and cpu cycles count which made the cpu jump to zero due to the offset being reset at the begining of a new execution cycle
| 
* commit 49c9d743dfdb77931608dab0aedd8ad54856e6d3
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Oct 7 22:15:15 2024 +0200
| 
|     adding the SUB instruction implementation
| 
* commit 53570ac5d6e8f7422ff5f543d7a90653978109e0
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Oct 7 18:08:34 2024 +0200
| 
|     when the debugger stop on a breakpoint or halt or stop after polling the new cpu state received, the debugger will also receive the ppu, apu & client memory state before stop the clock (please note that in this version, the crystal still doesn't stop exactly on the breakpoint or halt or stop if the crystal is running too fast meaning that a few ticks will already have been queued and will be processed as soon as the debugger finish receiving the current state)
| 
* commit 9c727138387f42f8ee6c0eb25d3cf344d32a0de4
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Oct 7 14:37:25 2024 +0200
| 
|     fixing instruction ldh [a8], A which was not correctly calculation the destination address which should be 0xFF00 + a8 // adding a test case for IR=0xE0 ldh [a8], A
| 
* commit 29703181dcab6bbac4ab42ea80a0220894e6a0c6
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Oct 7 09:10:24 2024 +0200
| 
|     fixed fetchOperandValue: ldh a, [a8] was not working because cpu wasn't getting value from operand [a8] correctly // added a corresponding TC in cpu_instructions_handlers_test
| 
* commit 88b10338bc31431db55e1fd9c0027ea95edee9c3
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Oct 5 20:03:16 2024 +0200
| 
|     launch listen loop on run and stop it when reaching a breakpoint or halt/stop instruction
| 
* commit 5bc6955a405b0c56303f508d16ca92b3f5566310
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Oct 5 08:25:14 2024 +0200
| 
|     gameboy should only tick ppu when lcdc register bit 7 is set to avoid race conditions on mmu writes
| 
* commit c92fb4d8aeb63158cae2682e14ad667b09eb0d7c
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Oct 4 23:51:04 2024 +0200
| 
|     state channel should transmit actual values and not pointers at least for memory writes since by the time the server reads the value, the mmu is already trying to access memory writes causing race conditions and the server to crash
| 
* commit d87643063fdbaa0b41548ee4cfa5b220ce903b34
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Oct 4 19:20:18 2024 +0200
| 
|     moving memory write to its own go file // adding tests for memory write custom marshal func
| 
* commit 78da089a2ea29630d21e716b72cac69f48f239ba
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Oct 4 15:12:55 2024 +0200
| 
|     memoryWrites custom json marshalling method was incorrect
| 
* commit 9a08d1ce118f6f9417ff9b917dd011106f5e5990
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Oct 4 13:43:37 2024 +0200
| 
|     gameboy.onTick: adding a busy channel to prevent ticks to imbricate
| 
* commit 5e40436ac99853288fb7fdd3415f1e00da9df015
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Oct 4 13:36:36 2024 +0200
| 
|     cpu onTick: adding busychannel to make the onTick go routines wait when then cpu is busy
| 
* commit abb30b01de4609bcf740b65063487585ddb73508
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Oct 4 00:13:04 2024 +0200
| 
|     clear mmu.memoryWrites in between steps // gameboy: increase ticks after the execution of the instruction
| 
* commit c76175be9d4b0bfb48ba50cc22d0acac1c2ee621
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Oct 3 23:41:38 2024 +0200
| 
|     v0.4.5: better channel management (optional channel, internal debugger state channels)
| 
* commit ed1c7a94465dd8871a62b21607d2ec25cbef3c3a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Oct 2 20:13:55 2024 +0200
| 
|     adding a test case for gameboy_test to check that all state channels are returning a new state when the gameboy is stepped
| 
* commit ec39287565e7d522073b0bcc903d3efa10e67505
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Oct 2 20:12:59 2024 +0200
| 
|     using mmu.getMemoryWrites in gameboy.onTick to report the memorywrites
| 
* commit bc15fa9c9ab8c509a4703e9fff2e8e70a1d8d3e4
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Oct 2 20:11:52 2024 +0200
| 
|     adding the getMemoryWrites receiver to mmu
| 
* commit b3201c73ed56b4cccfc030561abc3fad045d37fa
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Oct 2 20:11:17 2024 +0200
| 
|     simplify fifo.push by simply using pop to assure the struct keeps the max_length
| 
* commit 8c1ab3e7e0d3519366170df768f721182d50f0b1
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Oct 2 20:10:17 2024 +0200
| 
|     fixing a test case after cpu state interface changes
| 
* commit f2d817632cce8344c8b79671dac0ce5cd5894add
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Oct 2 20:09:17 2024 +0200
| 
|     adding test cases for debugger // found an issue with clientMemoryStateChannel not being set correctly causing the gameboy to stall on the server because it was waiting for memoryWrite message which never came
| 
* commit f0930bc4f107befb75bf27ad5a904a6e6456c86f
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Oct 2 18:25:47 2024 +0200
| 
|     updating the cpu, ppu and apu state json interfaces
| 
* commit 721816dc61665c59a486448d749465c34bf8b3e2
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Oct 2 15:17:52 2024 +0200
| 
|     handle all channel communication: cpu, ppu, apu & memorywrites from gameboy -> client and joystick state channel from client -> gameboy
| 
* commit 3389f3485d315cf10e20046874f3a7998716963e
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Oct 2 15:15:41 2024 +0200
| 
|     fifo: new first in first out data structure declared to use to hold the 100 last states transmitted by the gameboy state channels (cpu, ppu, apu, ...)
| 
* commit aeece0764b8fac07250b51c70cdaedf32fd73abd
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Oct 2 15:14:19 2024 +0200
| 
|     cpu instructions handlers test: exclude cpu.INSTR from compareCpuState since it won't add any value
| 
* commit 65ad86e744a3a31ed668a7062071ebb332bed7f8
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Oct 2 15:12:43 2024 +0200
| 
|     adding the instruction to the cpu and cpu_state
| 
* commit 51221fa630024e91f4379b9d4d38185bf15101a0
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Oct 2 15:09:06 2024 +0200
| 
|     apu & joypad: adding some dummy implementation to test the end to end communication between the front-end and gameboy-go
| 
* commit 9168f31a14ab0703a9eb5f406a76aab203eb9883
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Sep 30 13:21:44 2024 +0200
| 
|     cpu.onTick shouldn't increment the cpuCycles since this is currently done inside instructions
| 
* commit 7cc36765715587ca3ed9911e4514033fee9cdec4
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 29 20:20:22 2024 +0200
| 
|     main.go: forgot to align the NewDebugger call to its new interface which accepts cpu & ppu state channels
| 
* commit 6e7575e9effd0ad6c48c6dadf22f2fcc2621bdcd
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 29 20:16:48 2024 +0200
| 
|     align test cases with last changes
| 
* commit 7a2977fcc5fd9457ed963a45dedee8abd49f9929
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 29 20:15:58 2024 +0200
| 
|     make cpu & ppu state channels optional when instantiating a gameboy with NewGameboy
| 
* commit 4ffccdb7031e2bdde755ab66e2f66c33f15b650a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 29 20:13:48 2024 +0200
| 
|     removing logs to console to improve the gameboy performances
| 
* commit 96bd1888fcf5f58fec2204802d9d27442e5142a1
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 29 20:09:12 2024 +0200
| 
|     added an md file about the fmt.Sprint func showing that displaying logs to the console might harm the gameboy performances
| 
* commit 8a1f871c8b2df24152ea5d0faedaa9aa5bb4d8e7
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 29 12:59:49 2024 +0200
| 
|     adapted debugger calls to gameboy to new ppu channel
| 
* commit d83e855a8ba7ec30772e87c1ab6a147a6cd1463a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 29 12:59:25 2024 +0200
| 
|     implemented ppu state channel in gameboy
| 
* commit f201bf4fc8856370dee171d5a05eeea5e42fe21a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 29 12:58:51 2024 +0200
| 
|     implemented ppu mode switching and some dummy processing // added ppu state definition and modifiers
| 
* commit 81ea6820a20a2fbf03a92607c2ee6f3af9896f0f
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 29 12:57:10 2024 +0200
| 
|     Timer.tick should call subscriber.onTick in a separate go routine since its sole goal is to signal a new clock cycle to subscribers and not to manage any return value from the subs
| 
* commit 16d04af3d7dfd6206082d903b5aed7563bfb7da0
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 29 09:14:24 2024 +0200
| 
|     v0.4.0: gameboy-go package now reports on cpuStateChannel its cpu state with the outside world
| 
* commit 85da6d6372092d07eaaa6d95fbf28a54b713fe07
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 29 00:07:34 2024 +0200
| 
|     adding the cpu cycle value inside the cpu state
| 
* commit 70ade93dc5d320b16e2ebc13c49e8042f6d8cdf5
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Sep 28 16:06:51 2024 +0200
| 
|     fixed an error in a test: was waiting for the doneChan to be nil when it is instantiated with a bool chan in NewTimer
| 
* commit 1b26979db2e41c366be36941ecbc35e147e46d66
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Sep 28 16:05:14 2024 +0200
| 
|     skip not implemented tests instead of erroring to clarify what works vs what needs to be developed in CI/CD deploy pipeline unit testing stage
| 
* commit 306651a86e85a8695884b7058b1e533bb51570b3
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Sep 28 15:59:19 2024 +0200
| 
|     running only TDD unit tests and no BDD tests
|   
*   commit 05a5108878d97abdd59735fdf158c22cd809715d
|\  Merge: ed13649 2bea7bb
| | Author: codeFrite <codefrite@gmail.com>
| | Date:   Sat Sep 28 15:55:14 2024 +0200
| | 
| |     Merge branch 'main' of https://github.com/CodeFrite/gameboy-go
| | 
| * commit 2bea7bbd5585d1166bf6bd37685c7f02768e019c
| | Author: codefrite <34804976+CodeFrite@users.noreply.github.com>
| | Date:   Sat Sep 28 15:52:10 2024 +0200
| | 
| |     Update go version to 1.23.1
| | 
| * commit 7342fdf0bbedb927b9129afb4d2fe006d8ef4cdc
| | Author: codefrite <34804976+CodeFrite@users.noreply.github.com>
| | Date:   Sat Sep 28 15:50:32 2024 +0200
| | 
| |     unit-testing@main
| |     
| |     run unit tests on push to main
| | 
* | commit ed1364929cb135573ca4c6ebb6b9bceeb407e9c7
|/  Author: codeFrite <codefrite@gmail.com>
|   Date:   Sat Sep 28 15:54:46 2024 +0200
|   
|       fixed an error in main.go: NewDebugger requires 1 parameter for the cpu state channel notifications
| 
* commit 4dc902b968cad98338d19fed51f254e1c3bf1edf
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Sep 28 15:25:11 2024 +0200
| 
|     adding test cases for timer // fixing a bug in the Timer.Start that made the doneChan not working and the tick not triggering
| 
* commit bc9a59ea2a015c2c6a9cc127ef212776cb35319b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Sep 28 15:23:58 2024 +0200
| 
|     adding cpuStateChannel to Gameboy and Debugger
| 
* commit 49cdbab85fb9436d827aad6c7bbad74b51468115
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Sep 28 11:44:52 2024 +0200
| 
|     ignoring .gitignore
| 
* commit 32e0571ec44ff5a60f0cdd02104a11dc232421d9
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 27 23:15:38 2024 +0200
| 
|     move all unit tests to gameboy package to test private members of defined structs
| 
* commit 40edeb71f17a65313fdd73c4c053c22f1e901bfe
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 27 21:44:45 2024 +0200
| 
|     debugger state: changed the receiver from *Gameboy -> *Debugger and adapted all the calls // move print func from Gameboy.go -> Debugger.go
| 
* commit f11c7189b7fbfbe31ec79628a130924321c83726
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 27 15:20:37 2024 +0200
| 
|     renamed gameboy_state.go -> gameboy/debugger_state.go as it makes more sense since the cpu, ppu, apu members should all be private and only exposed through state report from gameboy or debugger info/interfaces
| 
* commit bc6922859a52e8397790b8d6ce2a45dc3f57e0fd
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 27 15:18:52 2024 +0200
| 
|     made all cpu members private : no need for them to be exposed outside the package: if needed user will use accessor from the debugger (might need to add some accessors from gameboy if needed to run games without debug info but still need something like cpu cyles
| 
* commit c2f1ace3f2e7c839822218dff72e680ea4e0979c
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 27 14:38:22 2024 +0200
| 
|     forgot some updates of the register access after changing them back to uint8/16
| 
* commit 914c0ac01bd4743ecb61f7b36e4590abf59901bc
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 27 14:35:54 2024 +0200
| 
|     added some docs about particular subjects: this way i'll get them up to date more easily
| 
* commit b3386dadd9884ba082c50bf8fddc7b29fcb79afb
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 27 14:34:52 2024 +0200
| 
|     renamed and adapted register8/16 unit tests
| 
* commit c596479ff30dfe3b0bb78efec974de96c2afa984
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 27 14:33:51 2024 +0200
| 
|     reverting back to using uint8/16 with helper funcs for registers that are NOT MAPPED to the memory mapping (easier to access and manipulate in calculation) // changing cpuCyclesCount from uint16 -> uint64 to have more room for counting frames (worth over 55y of emulator running time @60fps
| 
* commit facb9393ae867a99cb97e924d0987bc48f941117
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Sep 26 22:33:01 2024 +0200
| 
|     fix issue in gameboy_init: cartridge rom was created but not connected to the bus (should be able to attach more when implementing onboard SRAM or MBCs
| 
* commit bcd394756f58d5a58a6794527a58fae84227b7e6
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 19:02:11 2024 +0200
| 
|     fix Register16 implementation // add a test case for r8.ResetBit // add test cases for Register16 (similar to test for Register8 with larger numbers)
| 
* commit d35de51a2ec729b0022e75249f849cded8e4ae97
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 18:27:56 2024 +0200
| 
|     fixing a bug in register struct SetBit fct preventing the register to be assigned a new value
| 
* commit 2ca13233003edb12a37712cee89ce6177816c38e
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 18:18:25 2024 +0200
| 
|     adding unit tests for the Register8 struct
| 
* commit d5ce6761f4629d59b39dec11eb8fb76a34774e3a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 18:17:53 2024 +0200
| 
|     change package to test: please note that by default, even if the go.module points to our github repo and the test file imports github.com/codefrite/gameboy-go/gameboy, it actually uses the local version
| 
* commit 32c076d7224b3ca290b2292375414d087eb290cb
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 18:14:34 2024 +0200
| 
|     summarize the list of changes needed for the gameboy, cpu and ppu
| 
* commit 8c37258ff0f7cd99ae82462effe61e0167529208
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 18:13:33 2024 +0200
| 
|     adapt the references to the register now typed as Register8/16
| 
* commit 92e6557f81f5a9be7dc94a7109988124987ce48d
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 18:11:16 2024 +0200
| 
|     instead of uint8/uint16, use the new types Register8/16 to type the cpu and ppu registers and access some new features like getting/setting a particular bit inside a uint8/16
| 
* commit 40b342514f02e7869f833ee3a58cad1d8cb20178
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 17:54:52 2024 +0200
| 
|     add md file to prettier pipeline for table support
| 
* commit 6c577e1d2657bd0c799c8ff5fe209a2cc0b1ff4d
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 01:19:21 2024 +0200
| 
|     add a specific bus for the ppu 'ppuBus' and rename 'bus' to 'cpuBus': this will allow mmu to display different behaviours on data access through the specification of handlers to manage memory read/write: example: during ppu mode oam & lcd draw, oam and vram should be inaccessible by the cpu through the cpuBus but no restriction applies to the ppu
| 
* commit 252df6a4d3d8df3ebf7a5732f31e647de97768b8
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 01:06:12 2024 +0200
| 
|     ppu: adding NewCPU/onTick/updateLy methods
| 
* commit e3b849cd57d6a19a3fc22e5864bf79d0727e04a4
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 01:03:34 2024 +0200
| 
|     adding an interrupt handler on the cpu struct // define the clock tick handler to step the cpu
| 
* commit 626f0080f89b0fee578045728199ce25e25d7e51
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 01:01:30 2024 +0200
| 
|     set up cucumber godog project // define a feature with some scenario with their step definitions for the Timer struct (mostly used to test behavior, integration and end to end
| 
* commit bab4341caf3aba576960264aa3577533faa401c4
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 00:58:39 2024 +0200
| 
|     adding some comments
| 
* commit d3d0b4e88d8a5f63836e99bedac91ecc5f13b255
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 00:51:29 2024 +0200
| 
|     moving unit tests to 'test/tdd'
| 
* commit e048c5f87aa37d457a63a0e3f8ff71ff9d1ffa1d
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 00:48:31 2024 +0200
| 
|     adding the timer logic using channels
| 
* commit d7b4d6b39bdf869db01749826f06a4f5b9e8ab71
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 25 00:47:28 2024 +0200
| 
|     move all md files except for readme.md to a new folder 'docs'
| 
* commit e52c7f7c77ea88e778d99ea70a7ff18a349e25fe
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 20 10:39:15 2024 +0200
| 
|     start developing the timer struct
| 
* commit 58c55562a328721d012fecf6a523f0956aac94db
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Sep 19 05:03:17 2024 +0200
| 
|     INC instruction: flags were not set correctly at all
| 
* commit f3104538eb74423e47a5d47b8ea3fbfec5619c0b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Sep 19 05:00:02 2024 +0200
| 
|     the debugger now compares the states of 2 consecutive steps against the instruction json definition to see if anything that shouldn't has changed (flags, registers)
| 
* commit 7f6b363d49cc1dfb7fb578a9e3ea8268f8f210ed
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Sep 19 00:02:50 2024 +0200
| 
|     Need a stable version to develop a new feature in the front-end: the execution flow that shows the calls and frequences of use of instructions in the game program
| 
* commit fff8889782b75ea5c54671b67a4ff5f4d55cf876
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 18 02:03:11 2024 +0200
| 
|     aligned CALL & JP way of handling the offset calculation like in JR
| 
* commit f789d3648654ddfe621e0169fe962d9b2b7ecfbc
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 18 02:02:10 2024 +0200
| 
|     cpu.fetchOperand: fixed a bug in case 'a8' where the address was incorrectly read from memory
| 
* commit 9402ffd0ec766eda06bab31443cc61eb128c06fc
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Sep 17 23:26:26 2024 +0200
| 
|     JR instruction: when the condition is not met in Z/NZ/C/NC cases, we should still advance the pc to the next instruction start
| 
* commit 1186a81961e88655d4375a709552379bd4961a2d
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Sep 17 23:02:04 2024 +0200
| 
|     removed debug outputs
| 
* commit d442ef0df608cf437a2f04d7058113a5ad778edb
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Sep 17 22:04:57 2024 +0200
| 
|     modify debug output when running the gameboy : add more readable output at each step to see the chain of PC encountered during execution
| 
* commit 19e2c0540c8fd7f81724bda7586a7a0bb6db19a4
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Sep 17 22:03:04 2024 +0200
| 
|     update and fix test cases for the non cb prefixed instructions
| 
* commit 2d37dbcd9c677c07788c2e0bdc2c97e85ae7ee26
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Sep 17 22:02:23 2024 +0200
| 
|     adding the field CpuCycles to the CPU // make it appear in the cpu state // change cpu execution logic : add the update of the cycle in each instruction & update offset is done in instructions instead of the step loop
| 
* commit fd403b6b54913ffc9a16d733e8194e81c520f4ce
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Sep 17 08:52:30 2024 +0200
| 
|     cpu.fetchOperand: remove unnecessary console output that was being displayed every instruction and was polluting the output
| 
* commit 46b0a26807084e23871acfd24baedf28dc7e98b2
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Sep 17 08:50:50 2024 +0200
| 
|     updated go version requirement to latest v1.23.1
| 
* commit 17d4c411e68aa6cab488e5f989d99051bccef0a7
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 15 21:41:14 2024 +0200
| 
|     RET instruction was missing 4 cases: ret z/nz/c/nc / adpated the unit tests / added the operand z, nz & nc (c is taken by register) to retrieve operand values z, nz & nc
| 
* commit 77e0cb5e57abb07313a7ec99086f5da69c31ca02
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 15 16:46:23 2024 +0200
| 
|     implement reti instruction // add unit test for reti
| 
* commit 6e8bea42579c7a37297fc47f71896797da4b1b7c
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 15 14:43:02 2024 +0200
| 
|     added a unit test for the RET instruction
| 
* commit 7238c2ca1c866e28b8d3ec24697b950384f983f9
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 15 14:42:11 2024 +0200
| 
|     corrected a small typo in 2 comments
| 
* commit 8686e656284600147e68696766e66f34da12ecc3
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 15 13:21:34 2024 +0200
| 
|     added a unit test for the instruction JR // corrected some bad copy/paste in the labels
| 
* commit 1368c94c1c672de5032351200a665c0b50a8b28e
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 15 13:20:23 2024 +0200
| 
|     fix a bug found in unit testing: JR instruction shouldn't add the instruction length to the absolute address calculation and simply add PC with the e8 operand
| 
* commit 80f7715d70e1e36cd5ab5f6e9a0ec4a753227121
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 15 11:14:49 2024 +0200
| 
|     excluding .vscode from git version control
| 
* commit 690cec3c233be1c2ae0f90dabfeeadcdec6c9aa8
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 15 11:13:56 2024 +0200
| 
|     it has been in the git status output for to long :) whoosh ... it dissapeared ! ... I'll make the cpu work then i'll document the architecture at least (class diagrams, flow chart, use case, deployment)
| 
* commit 3c6de248de360d6c62878f98f6ac3549873fbf8f
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 15 11:12:20 2024 +0200
| 
|     excluding md files from prettier auto formatting
| 
* commit e42e891858097d596d5ad1da5806a640308d015e
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Sep 15 11:10:44 2024 +0200
| 
|     add unit tests for NOP, STOP, HALT, DI, EI, JP, CALL / corrected the instructions DI, EI, HALT, STOP by adding flags to the cpu to ask it to set/reset the IME after the execution of the next instruction so that we can finally pause the execution (breakpoint) on an instruction after a CALL/JP/JR opcodes / the step and run functions of the gameboy, cpu and debugger do now check if the cpu is halted or stopped before running anything, and if set, simply returns
| 
* commit fbfcdb3eb7f4ec46bf8349f0d54f4b9a7ca80605
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Sep 14 14:04:49 2024 +0200
| 
|     save the state of the gameboy after executing the test and before breaking out of the run loop to received from the front-end the last state after executing the current instruction pointed by the PC
| 
* commit 1e9cc5d94fbbdbf376079558f7477eea8982b6f9
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Sep 14 13:32:22 2024 +0200
| 
|     adding tc list header documentation // added comments on all tc to explain its purpose
| 
* commit 204aec813247e5a4d2026062aa6f726475f5e724
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Sep 14 11:15:36 2024 +0200
| 
|     during run, shift the state at each step to enable advanced debugger features like drawing the execution graph to highlight subroutines definition and construct the call tree
| 
* commit 4988fcc21abbc730310347e24ef6a40470134fb6
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Sep 14 11:14:13 2024 +0200
| 
|     enhancing the debugger output: while running, will print the status at each step and display the current state being executed
| 
* commit ab35c1bf4511f5bbdc74c81fbf1509d2d9c870b0
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Sep 14 11:10:55 2024 +0200
| 
|     randomize the cpu register states on startup to simulate the fact that registers are in an unknown state at startup
| 
* commit ec6aca982e880208c458ff8ecd0e95de569ee006
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Sep 14 00:03:18 2024 +0200
| 
|     adding blog posts to the blog.md files to documented things i have learned during the coding of this emulator
| 
* commit 43e9e153e1b589a3e48f7cc02e6b964cc1f05d46
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Sep 14 00:01:55 2024 +0200
| 
|     adding default case to trigger a panic error when an operand mode is not supported for a particular instruction
| 
* commit c3e8aea3a5a1d8087815b646455215f1b074b12f
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 13 17:40:03 2024 +0200
| 
|     fixing a bug in the TestMemoryDumpt test
| 
* commit b0625a05ab489cc7395b83774c8a7e281bc78e5b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 13 17:19:19 2024 +0200
| 
|     mmu.Write now does return an error if the memory location target is not found (i.e. memory not mapped)
| 
* commit b2b66eb8e498a9f35f4827632bd6906629b4c96a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 13 17:16:32 2024 +0200
| 
|     output the CPU State to the console and show how registers were affected
| 
* commit e797c0b0efccda5977a1879bf0437bd4efe69c6d
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 13 17:07:39 2024 +0200
| 
|     when asking the current state, the console will show the state of the memories mapped to the bus (name size : @from_address -> to_address
| 
* commit b4eb90bf07b24aa752933576a8c4fba794a15792
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Sep 11 12:40:05 2024 +0200
| 
|     mmu: since memory writes are solely used by the front-end, they should record the relative address of the change inside a particular memory instead of the absolute address to ease the front-end work (determine which memory location to update in a memory)
| 
* commit 5d5c492d61bf0eaea03a06551b0bb934ea574a20
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Sep 10 23:57:35 2024 +0200
| 
|     memory.GetMemoryMaps should include memory.Size() since in go the slice excludes the upper bound
| 
* commit 35372e41e9cd6f207489aa235b5d78db5e054fcf
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Sep 10 23:55:08 2024 +0200
| 
|     mmu.GetMemoryMaps should include memory.Size() since in go the slice excludes the upper bound
| 
* commit daa38457268b36759df8ac3a513a7c584a8aad3b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Sep 10 22:18:23 2024 +0200
| 
|     adding a marshalling method for uint8 arrays (JSONableSlice)
| 
* commit ba64decf8586b8f991fe785e768f71e99af716ce
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Sep 9 22:12:43 2024 +0200
| 
|     fixing an error in GetMemoryMaps which was dumping the memory of each memory up to memory.size instead of memory.size-1 causing the program to panic
| 
* commit 996ef1e326f14a9babf05d68ce1f1a6ccf69edf4
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Sep 9 21:52:22 2024 +0200
| 
|     the mmu should return the list of memories as memory writes and not memory maps as we do only care about the data attached to the struct and not the func that go with it // only 1 memory write should corresponding to a blob writing operation (we do not need to list all the different uint16 addresses impacted but rather the write begining address along with the data written to that address) ensuring lighter data exchanges with the server and more readable operation history
| 
* commit 7f391d5744ba911809df17c2a411e45a36787883
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Sep 9 11:49:55 2024 +0200
| 
|     add a getter to the memoryMaps to enable the front-end (debugger consumer) to know the state of the gameboy memories on startup
| 
* commit bcb6cac69fd21acfa847bc82b9fe56f73207728b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Sep 9 11:48:08 2024 +0200
| 
|     rename mmu.router field to 'memoryMaps' to differentiate it from memoryWrites // add a getter func to access memoryMaps for the debugger to know which memories are mapped to the mmu and access its data
| 
* commit efc1f27af44cd64c69aea8afde71cc7e0f3c3b64
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 6 22:27:31 2024 +0200
| 
|     when the package is not used as a dependency but directly used from the main, it still offers debugging capabilities: the main func instantiate a debugger and prompts the user for the next action to perform: add/remove/list breakpoints, run/step the gameboy, print the current cpu state
| 
* commit 3d31cb5ac1f1523c51c101a877fecfca80d91db3
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 6 22:24:25 2024 +0200
| 
|     add the debugger to allow the user to add breakpoints to stop the normal execution flow
| 
* commit 35b5b3106f4d05b26b9440653ee856a7dbb285dc
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 6 22:22:05 2024 +0200
| 
|     add MemoryWrite accessor and utility functions to print, clear and get the MemoryWrites contained in the mmu from the gameboy
| 
* commit e894b01ea3602f486cbfb443b11f8c467af37289
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 6 22:19:15 2024 +0200
| 
|     remove state operation (saveCurrentState, shiftState, ...) from gameboy as it is now handled by the debugger: this ensure us to have a fast running rom when run directly from gameboy and some debugging capabilities when run from the debugger
| 
* commit bb4af551d3fb89b76e63be7cdb1e8b51110dcec6
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 6 22:17:00 2024 +0200
| 
|     moving the MemoryWrite struct from gameboy_state to mmu // save all memory writes upon write & blob operations // adding a function to clear the memory write to allow the caller to clean it before a step or run
| 
* commit b4d7cca8fbd07f2e5513ff8423d20a10408bf378
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 6 22:10:30 2024 +0200
| 
|     minor change: specify which was successfully loaded before rendering it as a byte table
| 
* commit 79c02f5bc4e248ba0f087f7cdae218eac462acdc
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 6 22:09:11 2024 +0200
| 
|     adding doc strings to cpu.go
| 
* commit f850a150995c9e7e9ab64e5ec1516ab1f9106696
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Sep 6 22:06:54 2024 +0200
| 
|     adding doc strings to bus.go
| 
* commit 69e030bfbb89e93b31daa4b46930a6f812586fc6
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Jul 18 20:36:18 2024 +0200
| 
|     removing HRAM from cpu: it is now created by the gameboy itself // initializing the VRAM with random data to simulate the initial state after boot in a physical device
| 
* commit 22a9ec6894580bae288e8137266acf995b058668
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Jul 18 20:34:52 2024 +0200
| 
|     cartridge: using new func NewMemoryWithData in place of NewMemory + setData which was removed
| 
* commit ede3a092e52c77b3f5a860e6269cbdd975618d45
| Author: codeFrite <codefrite@gmail.com>
| Date:   Thu Jul 18 20:31:18 2024 +0200
| 
|     memory refactoring: removing setData func and replacing it with NewMemoryWithData to avoid memory changing its size // adding NewMemoryWithRandomData func to simulate the initial state of the VRAM on gameboy startup // adding tests to make the memory panic if we attempt to read, write or dump memory outside its size range
| 
* commit 20cb91583615826fe7e97e74b8e8ef259f51d06a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jul 16 22:11:45 2024 +0200
| 
|     moved tests in dedicated 'test' folder // adapted 'test.sh' script accordingly
| 
* commit 666e13392fdf3dc502097ef312ea550cc58f937c
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jul 16 22:08:59 2024 +0200
| 
|     instructions: removing unnecessary space character in increment and Decrement json tags in Operand struct
| 
* commit 1f82ec25e8456ea42f192a2c6a78b94c3f301ea3
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jul 14 13:25:06 2024 +0200
| 
|     BIT instruction was fetching the bit position to test from bus instead of getting it from the instruction (opcodes.json) which resulted in the test @0x000A JR NZ, e8(00FB) to never jump, clearing only the last bit of the RAM instead of clearing it completely
| 
* commit 37f27493303f12db9472928f4afac2c49a0cd1ab
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jul 14 12:45:39 2024 +0200
| 
|     INC HL was interpreted as INC[HL]: added immediate addressing mode for INC HL. Moreover, no flags should be impacted in that case
| 
* commit 810d56ed46d87c1368f18940c5eef3fdfe5a6982
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jul 14 11:59:50 2024 +0200
| 
|     now that JR X, e8 is fixed, i was stuck in an infinite loop between bootrom addresses 0x0098 and 0x00A1. The problem came from bad copy paste in DEC instruction that was also incorrectly setting the H flag for all DEC r8 instructions. This is now fixed
| 
* commit e482f1bb22180c758b3f1ab6a877cf9646d9a943
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jul 14 11:34:13 2024 +0200
| 
|     JR X, e8 apparently jumps relative to next instruction address and not current PC
| 
* commit 0e72be10cd7c4abaf7c4968a7ebee32be73c2940
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jul 14 11:13:13 2024 +0200
| 
|     fix opcode 0x20 JR XX, e8: the calculation to get the new pc was incorrect as it didn't take into consideration the fact that e8 is an unsigned 8 bit value making program counter backward jumps impossible
|   
*   commit f599873d41a0893aad9db0031a910f2c6a294741
|\  Merge: 0edd225 edede0c
| | Author: codeFrite <codefrite@gmail.com>
| | Date:   Thu Jul 11 00:06:32 2024 +0200
| | 
| |     merge after commit --amend
| | 
| * commit edede0c2a3d9caf17ebb26f6124601414ccceaf5
| | Author: codeFrite <codefrite@gmail.com>
| | Date:   Wed Jul 10 23:53:49 2024 +0200
| | 
| |     when instruction first operand is a memory location [a8] or [e8] we should fetch it manually in the instruction handler // updated LDH to fetch [a8]
| | 
* | commit 0edd2250b8a90f5c137bed09c9419f0d4bc4f40c
|/  Author: codeFrite <codefrite@gmail.com>
|   Date:   Wed Jul 10 23:53:49 2024 +0200
|   
|       when instruction first operand is a memory location [a8] or [e8] we should fetch it manually in the instruction handler // updated LDH to fetch [a8]
| 
* commit 24ad64942de4fb504ecebb653e111fe99316ca48
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Jul 10 23:44:47 2024 +0200
| 
|     [a8] operand are already given 'FF' prefix when return from fetchOperandValue
| 
* commit 69247dc38471c271b1bb2668250552d6e17ac71a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Jul 10 23:40:26 2024 +0200
| 
|     LDH instruction was using the wrong operand in case of LDH [a8], A for calcutating the address to which to copy register A
| 
* commit 77aa2617b3a260ee50c1bcf02b2434e7b4089543
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jul 9 23:12:25 2024 +0200
| 
|     fix XOR: registers N, H & C were not reset
| 
* commit 3e4e3c96fbe51562334efd84caf7319719010636
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Jul 8 11:06:07 2024 +0200
| 
|     fixed a bug that made the last bit of every memory unit unaccessible by the mmu, thus by the whole gameboy
| 
* commit d9ba5f9a3d40b96f418873414dd359c2c7381770
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Jul 8 08:19:10 2024 +0200
| 
|     fix io_registers and hram memories to be 0x0080 bytes and not 0x007F
| 
* commit 284c88d73872f8609327be8ec9c3f7b51c63b768
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jul 7 21:37:38 2024 +0200
| 
|     fix some offset values
| 
* commit 9bd3aae9937a950a5208b021135e37757900e5c7
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jul 7 21:31:23 2024 +0200
| 
|     fix some offset values
| 
* commit 9e8603677ec62b532e2c720bed0bab57f09d111a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jul 7 21:18:39 2024 +0200
| 
|     fix some offset values
| 
* commit 6965f4bf231487f150f09582617d01a187a1c474
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jul 7 21:13:30 2024 +0200
| 
|     cpu.offset behaviour change: its value is relative (cpu.pc=cpu.offset) and not anymore relative (cpu.pc+=cpu.offset): flow control instructions should set the offset to the absolute value and do not need anymore to make calculation to set the pc at the begining of the next execution cycle
| 
* commit f9a7256e23ce927a453e4fa12016074db02b4ab6
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jul 7 20:56:54 2024 +0200
| 
|     cpu_instructions_handlers: the offset moves the pc relatively to the current pc value (cpu.pc+=cpu.offset): fixing the code by assigning JUMP (ADDR - CPU.PC)
| 
* commit 532a5493955b471c542d36a32b254f2cdf555404
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jul 7 20:35:27 2024 +0200
| 
|     cpu_instructions_handlers: flow control instructions (call, jump, ...) will not set the program counter during their execution cycle: instead they set the offset which is treated after by the step func which begins by computing the correct PC ... this makes viewer work easier as data shown on screen are the result of the current instruction
| 
* commit 6a9025da99ca19ebf61670909ec9ec0f1d746b88
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jul 7 18:52:52 2024 +0200
| 
|     helper function currMemoryWrites which retrieves the whome memory accessible by the bus in the form of []MemoryWrites
| 
* commit 1abad021d872300bf8f8f41f3739b1ff6f5ebb5f
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jul 7 18:47:36 2024 +0200
| 
|     return the whole memory map instead of just the rom
| 
* commit 613edd7a02ba9b53c7002eb32a5b9bb30c1af642
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jul 7 18:47:01 2024 +0200
| 
|     adding field Name in struct MemoryWrite // add code accordingly
| 
* commit d7bd7b4b48890a504f9519bcea5f0aa1b5dce0ad
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jul 6 22:06:59 2024 +0200
| 
|     increment pc is performed before executing the next instruction to remove the lag between pc/instruction and the rest of the data transmitted as GameboyState
| 
* commit a97a83f0bdfa0e2af631f8585fbb395d9720e98b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jul 6 20:42:07 2024 +0200
| 
|     fixed the GameboyState.instruction being the last one
| 
* commit 1f6d6e08e5e545c246fa908dd5c55087bb117ef1
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jul 6 20:35:23 2024 +0200
| 
|     fixed the GameboyState.instruction being the last one
| 
* commit cd7becab0d54f122eb9aff2160b6c0fa9866d010
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jul 6 16:33:06 2024 +0200
| 
|     changing memorywrites.data type to []string
| 
* commit 58312884a28b92acbcc046e3b0c5849b10d22dea
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jul 6 16:28:07 2024 +0200
| 
|     integration testing: expose fixed code
| 
* commit 61bd6c9d5056149fcaef8a7fe344bc99e23ce28c
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jul 6 15:41:50 2024 +0200
| 
|     integration testing: check if code ok
| 
* commit 4811dfe26fd0cbbc65addd320946567545609723
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jul 6 13:18:39 2024 +0200
| 
|     fix once again the MemoryWrite struct
| 
* commit 6fc7651a1e235e9bde9a217c0d04d8bed3689346
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jul 6 13:00:56 2024 +0200
| 
|     fix MemoryWrite.Data signature to be able to correctly send memory bytes over web socket
| 
* commit 1f8c00aafb54eaa3b0f0d4fcc3649cc2729bebb9
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jul 6 12:09:59 2024 +0200
| 
|     GameboyState: modifyng MEMORY_WRITES type to []MemoryWrite
| 
* commit 83e0d42ce14bcddde96c0b045877a6cf9ad8f93a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jul 6 12:09:59 2024 +0200
| 
|     GameboyState: modifyng MEMORY_WRITES type to []MemoryWrite
|   
*   commit f62e764d4e0f7093021f7163283b9904c09de441
|\  Merge: 27258ee 8261527
| | Author: codeFrite <codefrite@gmail.com>
| | Date:   Sat Jul 6 10:38:52 2024 +0200
| | 
| |     Merge branch 'main' of https://github.com/CodeFrite/gameboy-go
| | 
| * commit 8261527cc5fc77ae75cffb79b0a4d043cb3d01ed
| | Author: codeFrite <codefrite@gmail.com>
| | Date:   Sat Jul 6 09:35:32 2024 +0200
| | 
| |     gameboy state: adding memory writes to the state
| | 
* | commit 27258eebdf908228b44913fd59124d8efabada61
|/  Author: codeFrite <codefrite@gmail.com>
|   Date:   Sat Jul 6 09:35:32 2024 +0200
|   
|       gameboy state: adding memory writes to the state
| 
* commit d7d856d3b014d09e5b0f24f75158904f25f889b9
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jul 6 09:32:47 2024 +0200
| 
|     go.mod: changing the module location from local to github
| 
* commit fa885e322d6e280d978f00ff8e98f17b62e4d9a9
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Jul 5 18:50:42 2024 +0200
| 
|     clean the local workspace
| 
* commit 5820a780b6f1c8e7db3fa7d213b679e9ea436d42
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Jul 5 18:13:50 2024 +0200
| 
|     gameboy: complete reworking of the public interface // gameboy initialization is done differently: it now requires the rom name and returns the gameboy initial state (all 0's, nil, ...)
| 
* commit be4d98ac13b34298fe8de5a414e81723cfdbb33e
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Jul 5 18:11:57 2024 +0200
| 
|     cpu: all members of CPU struct are now public to be accessible by JSON functions // fetchOperandValue and fetchOpCode do not update the cpu state anymore but rather return the value to the caller // minor cosmetic changes (comments and variables order)
| 
* commit 0fd70b3698eeb583fd7a185dfe58c89b7ecf48e1
| Author: codeFrite <codefrite@gmail.com>
| Date:   Fri Jul 5 18:07:49 2024 +0200
| 
|     go.mod: changing go version from 1.22.4 to 1.22
| 
* commit 3c4793993b78fe355f43af296d84a50b2da6e4be
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jul 2 23:08:22 2024 +0200
| 
|     get cpu state (registers & current instruction being executed) from the outside (gameboy struct)
| 
* commit 40f12a2e6a0085cb8e86f811e6e4499003a7f9ae
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jul 2 22:59:03 2024 +0200
| 
|     code automatically reindented after save
| 
* commit 7069e36c9067743a60b3debeb4a6eea215ab6846
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jul 2 22:16:55 2024 +0200
| 
|     utilities.PrintByteTable: add color to the terminal output and add a param to only show a certain number of 16 bytes lines
| 
* commit b2eb06640ced66f9de006e3821dc8bd81f1069f8
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jul 2 22:11:32 2024 +0200
| 
|     go.mod: changed the go version from 1.18 -> 1.22.4
| 
* commit 3b35d480238089a987ca3e0d25fab34ec1df0563
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jul 2 22:10:18 2024 +0200
| 
|     split gameboy/cpu.go into multiple files
| 
* commit 11f12e0cf28affccaf680a0f2b40a8e752c5cff8
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Jun 24 19:19:20 2024 +0200
| 
|     add implementation for POP instruction
| 
* commit dcc4d2ef4cd40b4b2da4e372b711b9f5d8ed1d89
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Jun 24 19:18:52 2024 +0200
| 
|     LD instruction: manage case LD [a16], r8/r16
| 
* commit 79204c594beca58d466ee2cc858dd4094993f2e1
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Jun 24 19:17:59 2024 +0200
| 
|     add implementation for instruction RST which jumps to a fixed vector location /bin/zsh, , , , , ,  or
| 
* commit 05cf3afdbfae6d1e097535760464bae5e27e9112
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Jun 24 19:15:53 2024 +0200
| 
|     fetchOperandValue: manage cases HL, [HL], SP and all fixed call vectors /bin/zsh, , , , , ,  &  used in the RST instruction
| 
* commit 12430cdff53c4f9d2ab445bb5b4d1460bf292fdf
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Jun 24 19:12:57 2024 +0200
| 
|     utilities: externalize the code snippet that prints the memory content as a table under PrintByteTable(data []byte)
| 
* commit 55979b29dffca447777a543b73d778d571af9c2a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 23:18:43 2024 +0200
| 
|     RL instruction implementation:  rotate r8 or [HL] left through carry: old bit 7 to Carry flag, new bit 0 to bit 7
| 
* commit 629355755bce8e6dd3512f8436d0ff1470c82df2
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 22:19:09 2024 +0200
| 
|     adding PUSH instruction implementation
| 
* commit 690bea531b90f01d13b882ec01e653bf6cb70f36
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 21:43:50 2024 +0200
| 
|     cpu: adding support for BC & DE registers in func fetchOperandValue
| 
* commit 8c83b76999cab837300300e21de9970bad0d4227
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 21:43:02 2024 +0200
| 
|     main: change the path to the rom
| 
* commit 922f8acc4e1ad92fb654c6661e441e89cc565be4
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 21:42:31 2024 +0200
| 
|     gameboy: map the cartridge on bus @0x0100 to leave the 0x0000-0x00FF available for the boot room on power up
| 
* commit 0badf7c2013cb65eec70e16d6279027c84a883fa
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 21:37:21 2024 +0200
| 
|     cartridge: replace rom []byte type to a pointer to a *Memory. Adapt the struct funcs implementation accordingly
| 
* commit b261b6655b33db60780c1eb534fa3d655ef984f8
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 21:25:54 2024 +0200
| 
|     memory: adding a setter setData(data []uint8) to set at once the value of the memory data
| 
* commit cfeebbd108828d06df5debb836b5b1b535040b1d
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 10:02:49 2024 +0200
| 
|     LDH instruction: should be adding 0xFF00 to [a8] when reading/writing from memory
| 
* commit b2dc714744d0caa852c70493c3d620a11cc184a1
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 10:00:52 2024 +0200
| 
|     cpu: adding Boot func to execute the boot rom and returning when reaching PC=0x0100
| 
* commit c3570df7d4e10637d9fb3e6019db4ea8f3562640
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 10:00:06 2024 +0200
| 
|     cpu.printRegisters: spliting logging output to 2 lines to make it more readable
| 
* commit 738e5bad7d2600d30717dcda5d05cd534901e50b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 09:31:53 2024 +0200
| 
|     cpu: all r16 getters were incorrect: i was shifting << 8 r8 before casting to uint16 which caused the high nibble to be 0x00. Fixed it for BC, DE & HL
| 
* commit 20d8f65dc84b74aed016dcb7220d0d949fcb1072
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 09:09:23 2024 +0200
| 
|     LD HL: post-incr/decr was already done after at the end of the case HL. I removed the duplicate
| 
* commit dab97ce0bef3ca2cfa46f5cb09d612993967b7d6
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 09:04:40 2024 +0200
| 
|     XOR: setting the Z flag properly after operation
| 
* commit d98c26bb826006d1c62e41d65d3ad0e7aa651a36
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 09:01:36 2024 +0200
| 
|     LD instruction: taking care of post incr/decr of non immediate [HL+/-] variable
| 
* commit 897ee8ec052c8883cd62c47339cf334bf859acf0
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 08:56:54 2024 +0200
| 
|     adding func cpu.printRegisters to display the content of the r8 registers (A,B,C,D,E,H,L) and the flags (Z,N,H,C)
| 
* commit e2caeb61cd17ea56ae7ad5dce5c4256d12d12e6a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 08:55:27 2024 +0200
| 
|     cpu.printCurrentInstruction: displaying correctly operands that are incremented/decremented by appending the + or - sign to them
| 
* commit cd6f0be4f721d48e692bbe6c03491782f0279f2e
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 08:52:57 2024 +0200
| 
|     adding HRAM @0xFF80 of size 0x007F on initialization
| 
* commit d6a6f908d21da8a8f650713504fb2d490d04e3e7
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sun Jun 23 00:06:55 2024 +0200
| 
|     cpu: replacing and deleteing the PREFIX instruction by a boolean inside the CPU struct and use it in the code to determine which 'execute instruction' to call
| 
* commit 353ec8b0cb44d388cbcc088200e4260682c13220
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jun 22 23:10:47 2024 +0200
| 
|     refactoring: removing complex logic from bus and moving it to mmu in order to be able to switch address 0x00 from boot rom to cartidge ROM (not yet there but works as before // inverting params in func AttachMemory: address comes first, then memory (interface Accessible) and update reference accross code (bus, gameboy)
| 
* commit 6d040eb5e6b4aa53af89519606a2272746900b88
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jun 22 23:07:07 2024 +0200
| 
|     bug fixing: increment PC when calling cb prefixed BIT instruction
| 
* commit 8034a1df1f08473c2cf53a82377b4388945088aa
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jun 22 23:05:42 2024 +0200
| 
|     adding implementation for not prefixed instruction INC
| 
* commit b005521a1261217c3540df43f55cdb550e1eb2f6
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jun 22 23:04:04 2024 +0200
| 
|     utilities: adding header (byte position) before writing  the content of the loaded rom
| 
* commit c39baf55006562fbac1ca79bae9468001d4919ba
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jun 22 23:01:42 2024 +0200
| 
|     adding printCurrentInstruction func that logs the current instruction being executed in the format  * PC: 0x00A6, Bytes: EA AB 01, ASM: LD [AB], HL and using it in the step function
| 
* commit 378cbf518ce8d9fd8d0e98efd8e4f056444cbae9
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jun 22 08:04:00 2024 +0200
| 
|     adding BIT instruction implementation
| 
* commit 89ed1e5f19efda76fc8c4ca947c6ee45fcf20c5a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jun 22 07:53:13 2024 +0200
| 
|     bus: adding a WriteBlob(addr uint16, blob []uint6) func to bulk write byte to the memory: used in gameboy.go to load the bootrom to memory
| 
* commit f5a5a58f55901b87d1816a9b359bdbe9a9e0661d
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jun 22 07:51:16 2024 +0200
| 
|     Moving the loadRom func from cartridge to utilies.go and make it public
| 
* commit 104c5d844e9ea4b9beb661d1292ea08f17021853
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jun 22 07:49:33 2024 +0200
| 
|     adding file gameboy/gameboy.go to initialize a gameboy CPU/BUS and its memory along with a memory to hold the boot rom
| 
* commit 3d1c23b6b74ce5402a3238891c23050bc7ef6378
| Author: codeFrite <codefrite@gmail.com>
| Date:   Sat Jun 22 07:48:04 2024 +0200
| 
|     rename gameboy.go -> main.go
| 
* commit e0a70629630ba372d6ea3d5bc89e3f888694a21a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Jun 19 16:08:12 2024 +0200
| 
|     add CPL instruction that compares two memory locations (after implementing this instruction, when I run a rom, I end up calling an illegal instruction. This can either come from a defect in my implementation OR an issue linked to the fact that my program is not booting from the boot rom which results in an incorrect starting state ==> I will first include the boot rom to the initialization phase and then rerun the ROM and see if the error dissapear. Otherwise, I'll be adding tests for every single instruction to make that the implementation is ok
| 
* commit b3ecdbffdbd84c59dfcea46fabb10d72d5318a4e
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Jun 19 15:53:06 2024 +0200
| 
|     adding LDH instruction as well as the corresponding memory location I/O registers from 0xFF00 to 0xFF7F
| 
* commit 680c2c0a4c71fd083dfc86b1e0d1654413116bd2
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Jun 19 11:13:18 2024 +0200
| 
|     cpu.fetchOperandValue: implement case 'e8' which saves a int8 (signed 8 bits integer) to the cpu.Operand context // add the opcode to the panic message when operand type (Operand.name) is unknown
| 
* commit 0690fd90c8a28de9eacf78826d7ce5445541dbb6
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Jun 19 11:09:49 2024 +0200
| 
|     LD: taking care of the 0xF8 instruction LD HL, SP+e8 where 3 operands are present and where flags need to be updated (all other LD instructions leave the flags untouched
| 
* commit 809c08242d02fae24514f62330b3789eaa4e0b5b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Jun 19 10:58:57 2024 +0200
| 
|     adding the JR instruction
| 
* commit f7860e754239ccd24361bcacc3777b40a5b9f377
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Jun 19 08:44:08 2024 +0200
| 
|     ADD instruction: fixed the address pushed onto the stack: i was pushing PC instead of the address of the next instruction which is PC + instruction.bytes
| 
* commit 22525779a41d9cd570172b3cbb815a6e9bc9bb1f
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Jun 19 08:30:37 2024 +0200
| 
|     adding DEC instruction
| 
* commit 00b71207f1475b90f6d2fd023bd20e38c58b46ce
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Jun 19 00:50:45 2024 +0200
| 
|     cpu_instructions_handlers.go: adding CALL, JP & LD implementations // refactoring RET instruction by using bus.write16 func (instead of writing manually low then high bytes to memory
| 
* commit 0bc06714c4078c96185590b596abb10bce7593db
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Jun 19 00:47:59 2024 +0200
| 
|     cpu: adding setters for DE & HL registers // adding stack operations push & pop // adding debug outputs in step func to see which instructions along with operands the cpu is currently executing
| 
* commit bac74cc5596f78b3963be45b4c674909fc846335
| Author: codeFrite <codefrite@gmail.com>
| Date:   Wed Jun 19 00:45:49 2024 +0200
| 
|     gameboy: add initCPU func to simulate the effect of the boot room which sets some registers and memory location to particular values (init not yet completed)
| 
* commit e5e84fcc1337aba14414147f93473d194ab5c47b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jun 18 22:39:24 2024 +0200
| 
|     instructions_test: adapt cpu.ExecuteInstruction call: no more operand1 & 2 params needed since operand is now saved to cpu context
| 
* commit ff5219bffa23ae27103442909e43ae4b02f01ab3
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jun 18 22:37:33 2024 +0200
| 
|     cpu: finalize the step() func: fetch last operand value (if any) and move c.executeInstruction call to the end to mimic fetch, decode, execute CISC cpu instruction cycle
| 
* commit 037e199b4074f9a5bf6ec961603786486218b870
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jun 18 22:31:48 2024 +0200
| 
|     cpu.fetchOperandValue: refactoring the func to decode the instruction operand value as an uint16 which will be used saved to the CPU.Operand context
| 
* commit b038178783b211c175abf4af01922cb0aa04b37e
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jun 18 22:29:32 2024 +0200
| 
|     bus: adding the func (b *Bus) Read16(addr uint16) uint16to read 2 consecutives bytes from the bus and return it as little-endian
| 
* commit b0d9d65b153be3e4f1c9883231317c05c0a86875
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jun 18 22:03:46 2024 +0200
| 
|     cpu: add Operand uint16 field to get track of the operand resulting from the decoding of instructions in order to not have to pass the param to all 512 instructions and to the executeInstruction & executeCBInstruction // cpu refactoring: reordering the different func to match the fetch, decode, execute CPU cycle
| 
* commit effcf99ffd78a66635700d1b72ed82450362dd8b
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jun 18 21:33:14 2024 +0200
| 
|     instructions: getting the caller path to properly load the opcodes.json file in both situation where i run the program with 'go run .' and when i run the tests 'go test -v ./gameboy'
| 
* commit 8cf51ef42900888a229d3adfa5c2e781ade915f5
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jun 18 21:20:13 2024 +0200
| 
|     .gitignore: added macos .DS_Store file to the exclusion list
| 
* commit cc0d8fadd238e166215a59f4ffa3eb13b7a2fdd0
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jun 18 21:18:55 2024 +0200
| 
|     cpu: moved non-prefixed and prefixed instructions handlers to dedicated go file for lisibility and to isolate the main cpu functions: run, step, fetchopcode, fetchoperand and various setters & getters
| 
* commit 92707401fc9dbb4862e94f1e9d59576c6ddb988a
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jun 18 12:02:25 2024 +0200
| 
|     added implementation for all illegal instructions 0xD3, 0xDB, 0xDD, 0xE3, 0xE4, 0xEB, 0xEC, 0xED, 0xF4, 0xFC, 0xFD routed from executeInstruction default handler to CPU.ILLEGAL() func which triggers a panic
| 
* commit ce7ea32f3ec2abdbd59078b4c6c6b225d1474245
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jun 18 11:59:36 2024 +0200
| 
|     reordered illegal instructions and numbered them correctly
| 
* commit c5bcbc677b2262f20bdadfd16c4bfefb4f5fc914
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jun 18 09:26:19 2024 +0200
| 
|     added test case for 0xCB PREFIX // added an utility func for catching panic when calling a func, allowing the test execution to continue: this function returns a boolean indicating if the call resulted in a panic or not
| 
* commit da0bf6c099a1162b8583c9d0fdb2a5e223ac239f
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jun 18 09:24:39 2024 +0200
| 
|     added instruction 0xCB PREFIX along with the CBInstruction router executeCBInstruction // reorganised CB instructions handlers to match their order in opcodes.json
| 
* commit bb083c12421a6ac6e306e67494ed8148ebdb9722
| Author: codeFrite <codefrite@gmail.com>
| Date:   Tue Jun 18 08:20:53 2024 +0200
| 
|     added 0xC9 RET instructions // test case added but waiting for PUSH instruction to implement (now skipping it)
| 
* commit 8e9c798459090983c852e5453b9191375dd8124c
| Author: codeFrite <codefrite@gmail.com>
| Date:   Mon Jun 17 22:03:20 2024 +0200
| 
|     HALT instruction: added halter flag to CPU struct // made the CPU step func private and introduced the public func to be called from main: calls step in a loop if halted in reset, otherwise break out of the game loop (should be changed in the future to wait for interrupts) // added corresponding test case (to be completed when interrupts are developed)
| 
* commit 3c85c7daf90b116a94b2db00e9732c17fb4aa673
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 17 14:57:52 2024 +0200
| 
|     adding code for 0 operand instructions RLCA, RRCA, RLA, RRA, DAA, CPL, SCF, CCF
| 
* commit 85a6e479eaeb1707b6b8389040ef510857c518fc
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 17 14:56:34 2024 +0200
| 
|     cpu.go: modifying the steps in the 'Step' func to better mimic the CISC instruction evaluation cycle
| 
* commit bfbda5926580193ca5efd1cc4ea11e791b65966f
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 17 14:51:58 2024 +0200
| 
|     changing the uri to the 'opcodes.json' file // making LoadJSONOpcodeTable() func public
| 
* commit 35f5598b7e2ebfbd5f1ab0c91bdc4bb563c21965
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 17 14:49:23 2024 +0200
| 
|     deleting utilities.go and related TCs
| 
* commit 3444d65866a5c93c2a2159e1daf7ea5dff8b8f5c
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 17 14:48:27 2024 +0200
| 
|     updating documentation
| 
* commit 0ca15a287e0ce1f318cd2330e465fcf7b294425e
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 17 14:47:29 2024 +0200
| 
|     splitting all uint16 registers to uint8 to mimic/match the gameboy cpu 8-bits architecture // adding get/set/reset/toggle func for flags // adding getter/setter for accessing uint8 registers as uint16 (BC, DE, HL)
| 
* commit 3ab525348aa7beb824d50b72438c1f39b2c8ef25
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 17 14:39:40 2024 +0200
| 
|     gameboy.go (main): refactoring to use the new memory struct instead of vram & wram
| 
* commit ae45e25dfaad0aa4286bba2bee451c6e7bd692e2
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 17 14:38:13 2024 +0200
| 
|     refactoring the bus to use  the new memory struct instead of vram & wram
| 
* commit 9c2a3609025613b1a47e07eadcdacdb71480501a
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 17 14:37:07 2024 +0200
| 
|     refactoring code to use the new memory struct instead of vram & wram // modifying all []byte types to uint8 // adding Size() func to get the rom size in bytes
| 
* commit 9338e4ae0b3b91ba83553d11a71c0788f95b5116
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 17 14:15:58 2024 +0200
| 
|     replaces vram.go & wram.go with memory.go // renaming and redactoring memory_test
| 
* commit 86ec30dfd6fc2cb3bd0bdc634df4b8e9dbe25245
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 17 14:13:18 2024 +0200
| 
|     test.sh: adding yellow formatting for skipped TCs
| 
* commit 8ddb4f76ea28c6ea4577801c1c03c666c1affe6d
| Author: codeFrite <your_email@example.com>
| Date:   Wed Jun 12 10:34:54 2024 +0200
| 
|     adding optionals fields increment & decrement to the operand struct
| 
* commit cb7ee0170d879766a483f26ca5324682d69541fc
| Author: codeFrite <your_email@example.com>
| Date:   Wed Jun 12 10:30:46 2024 +0200
| 
|     changing completely the approach: opcodes are now mapped to CPU receivers, instructions are dealt with in batch (all LD, all JP, ...) and operands are read from the instruction json definition and not individually in each opcode
| 
* commit 5840082b451937b205723a683dfb2bf052690a88
| Author: codeFrite <your_email@example.com>
| Date:   Wed Jun 12 10:29:04 2024 +0200
| 
|     set PC to 0x0100 to skip bootrom and access ROM game begining // remove unused var i since the program stop by itself by returning a panic if an unknown opcode is found
| 
* commit 251f259e4cf4e67c0837099f6fd2c9fc62f4d0da
| Author: codeFrite <your_email@example.com>
| Date:   Tue Jun 11 08:37:36 2024 +0200
| 
|     adding opcodes from https://gbdev.io/gb-opcodes/Opcodes.json in order to code them more efficiently and use a community accepted 'source of thruth'
| 
* commit ca9af11ef1ea253af8340e767be22059639bff1c
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 10 21:29:18 2024 +0200
| 
|     adding test case for 0x4A LD_C_D
| 
* commit d033efdd54f2f187b09be3f897c222ee31527d46
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 10 21:07:28 2024 +0200
| 
|     opcode 0x2C INC_L: adding flags update logic and corresponding edge TCs
| 
* commit 4700be08b9f7be3f563f05b4c46fa9874e5e501b
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 10 21:06:07 2024 +0200
| 
|     adding getter and setters for F register flags Z, N & H
| 
* commit ab18c4e7fee2ef58763d65152cf81b4a5dd7ce53
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 10 12:51:10 2024 +0200
| 
|     adding new test case for instruction 0x2C INC_L which increments the value of the L register
| 
* commit d52be24af128c972a9d8e293b48b9fd05958e588
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 10 12:50:19 2024 +0200
| 
|     fix instruction 0x11 LD_DE_n16 which was loading the next 2 bytes in the wrong order inside DE register
| 
* commit 3d9f1b9fb85e070ea6987c469aba0c3f6e806ed5
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 10 11:46:32 2024 +0200
| 
|     adding some unit tests for the instructions func (0x00, 0x01, 0x11): uses WRAM to load the instruction and operands to be executed by the CPU
| 
* commit 15370e149db2b4c55242351e2e71c4eba885c4e6
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 10 11:44:51 2024 +0200
| 
|     adding implementation for opcodes 0x11, 0x2C, 0x4A, 0x4B, 0x53
| 
* commit d1b36affb034b07c70d39e37d610c818960b7c9a
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 10 11:43:48 2024 +0200
| 
|     adding a test script that launches go text and add some formatting (GREEN for passed tests and RED for failed ones)
| 
* commit e23c1d590380a60a3687bd73bdd4afd631353c9e
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 10 11:43:07 2024 +0200
| 
|     adding test cases for utilities functions uint16 to bytes and vice versa
| 
* commit e904b41fd28924e0a8b5884c81ba250cdc9538ab
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 10 11:42:19 2024 +0200
| 
|     adding write func to the bus. Writes a byte to the cartridge, VRAM or WRAM depending on the address: func (b *Bus) Write(addr [2]byte, value byte)
| 
* commit 0aa20b13b1797c63e50801245c4304cd5aa8ab05
| Author: codeFrite <your_email@example.com>
| Date:   Mon Jun 10 11:40:47 2024 +0200
| 
|     CPU returns an error if instruction has not been yet implemented or is unrecognized
| 
* commit 7f4632d7d31c2567d5967aa2c751793096c6b564
  Author: codeFrite <your_email@example.com>
  Date:   Sat Jun 8 23:46:19 2024 +0200
  
      initial commit // rom: load the binary file and retrieves the main header fields // bus: allows access to rom, vram, wram // cpu: main logic for fetching opcode and operand, executing the instruction, updating the internal registers, displaying the cpu register content (very limited set of opcode supported)
