# dalek-manager
A web config tool to be hosted on the roborio with several team code infrastructure diagnostics

## File Structure
json will be used for most things

### Binary Directory
The robot code binary is located at `/home/lvuser/`.
 - `FRCUserProgram` - Robot Code binary run at startup by the roborio
 - `dman/` - Working directory for the application

### Working Directory
The working directory is `/home/lvuser/dman/`. Within it is the manifest file and several subdirectories
- `manifest.json` - This is a file containing file paths to various files.
- `autonomous/` - This directory contains all autonomous scripts
- `ports/` - This directory contains all port configs
- `controls/` - This directory contains all controller bind profiles
- `settings/` - This directory contains all user-defined config files
- `logs/` - This directory contains all previous logs and a ledger of all logs
- `binaries/` - This directory contains all versions of the robot code binary and a ledger/table of info on each one

### Manifest
The manifest file contains the following attributes
 - `templates` - Contains file paths to template files
  - `configs` - Contains template file paths for configs (ports/controls/settings)
    - `ports` - Path to ports template file
    - `controls`
      - `available` - Path to available controls file
      - `requirements` - Path to required controls file
   - `settings` - Path to settings template file
   - `log-config` - Path to log config template file
 - `runtime` - Contains file paths to files loaded by robot code at start
  - `configs` - Path to configs
    - `ports` - Path to the json file containing component port labels
    - `controls` - Path to the json file containing the controller binds profile
    - `settings` - Path to the json file containing end-user-defined settings for the robot configuration
    - `log-config` - Path to the json file containing log configuration
  - `autonomous` - Path to the lua script for autonomous mode
  - `binary` - Path to the binary being run
 - `server` - Settings related to the app
  - `web-root` - Path to web root
  - `port` - Port to listen on

### Configs
Configuration files encompass controls, ports, settings, and anything miscellaneous that fit into this category. They are stored as json files.

#### Templates
These are files generated by the robot code which specify the data which need to be included in the configuration file. These are read by the manager app in order to generate the web interface allowing the user to edit the config files. For example, in ports/template.json is a layout of the names of hardware components on the robot as well as the range of port numbers or IDs which they can be mapped to.

##### Ports
*The following explains the structure of the ports/template.json file generated by the robot code.*
The file is an array of "port spaces". A port space is essentially a type of input that has its own set of ports, e.g. CAN Bus IDs are enumerated independently from solenoid outputs or digitial input-output pins.
 - Port Space Object
  - `name`: The name of the port space `"[CAN|DIO|solenoids|pwm|analog-in]"`
  - `min`: Minimum port or id used in this space
  - `max`: Maximum port or id used in this space
  - `keys`: An array of either strings which correspond to hardware components that have an independent port or objects which correspond to a named subset of the port-space, simply organizational.

##### Settings
*The following explains the structure of the settings/template.json file generated by the robot code.*
The file contains an object with a single member
 - `subsystems` : Array of conceptual subsystems containing settings which pertain to them
  - `name` : Name of the subsystem or setting
  - `values` : Array of settings or subsystems
   - `type` : Specifies the type of data the setting is. If this is null or undefined then this object is a subsystem. `"[integer|float|boolean|string|null]"`

##### Required Controls
*The following explains the structure of the controls/requirements.json file generated by the robot code.*
The file contains data which tells the configurer the commands and systems on the robot which require a teleoperated input to be bound to them.
At the root is an array of subsystem objects.
 - Subsystem object
  - `name` : Name of the subsystem
  - `continuous-inputs` : Array of continuous action objects. These are robot components which require updated input every frame, such as the drive base and most other analog controls.
     - `name` : Name of the action requiring control
     - `type` : Type of parameter needed `"[analog|digital]"`
  - `commands` : Array of command objects. These are processes which can be bound to a trigger.
     - `name` : Name of the command

##### Available Controls
*The following explains the structure of the controls/available.json file generated by the robot code.*
The file contains data which tells the configurer which inputs are available to be bound to commands and such on the robot.
At the root is an array of input space objects
 - Input space object
  - `name` : Name of the input space object (e.g. "joysticks")
  - `ids` : The presence of this object means that this input space has several of the same type of input which are ID'd by the scheme described in the object. One example of this would be the multiple joysticks.
   - `min` : Minimum ID
   - `max` : Maximum ID
   - `names` : Array of strings corresponding to possible named ids
  - `types` : An array of input space objects. The presence of this array means that there are multiple types of inputs on the same ID space. e.g. joysticks and gamepads.
  - `digitals` : An array of strings corresponding to various "digital" inputs (swithces or buttons)
  - `analogs` : An array of strings corresponding to various "analog" inputs (joysticks axes or analog sensors)

#### Output Configs
The web config is responsible for outputting a json file which is read back into the robot code to tell it how things are bound or configured. The robot code also generates a default.json along with the templates so that it can read in default behavior as well as provide the web config with a default value.

##### Controls
*The following explains the structure of the ports/default.json file generated by the robot code.*
This file serves as a model for other output control configs as well. At the root is an object with two members
 - `binds` : Object which contains subsystem objects, named by `controls/requirements.json`
  - Each subsystem object a series of commands or action objects, known to the robot by the key
    - Each command or action object contains info on the nature of the bind
      - `source` : Either an object or id of an object which corresponds to an input source
        - `input` : Address of the input. Each input space is separated by a slash. The available inputs are defined in `controls/available.json`. For input sources which have types as well as IDs, the type is referenced first, and then the ID. For example: `joysticks/joystick/0/y-axis`
       - `sources` : An array of sub-source objects. Its manipulation depends on the other variables
       - `source` : A sub-source object. Its manipulation depends on the other variables
       - `type` : Defines the type of source
         - `null` : If the type is null it is just a direct input source
         - `toggle` : This requires there to be a boolean `source` object or input in the source. This specifies that the value of the source flips every time the source switches from true to false
         - `and` : This requires there to be an underlying `sources` array in the source. This specifies that the value of the source is true if all of the elements in the `sources` array are also true
         - `or` : Same as `and` but is true if any of the underlying sources are true
         - `none-of` : Same as `and` but is true if none of the underlying sources are true
         - `compare` : This requires there to be a `comparison` member as well as either a `source` or `input`.
        - `comparator` : Value of the same type as the input which is compared to the input
        - `comparison` : If the source is of type `compare` this specifies how the input is compared. It supports: `">", "<", ">=", "<=", "==", and "!="`
       - `poll`: This specifies the event which leads to the command being run
         - `onTrue` : Runs the command when the source becomes true
         - `onFalse` : Runs the command when the source becomes false
         - `onSwitch` : Runs the command when the source changes
 
### Subdirectories
Most of the subdirectories are simply just a collection of various versions of the file associated.
Most also have a `template.json` which is generated by the robot code client which contains a template of the values and defaults requested by the robot code at initialization. This is used by the web app to assist in initializing values in a new or outdated config.

#### Logs
The `logs/` directory contains several more directories indexed by run time. Each will contain a copy of the `manifest.json` used for that runtime at startup. Each log directory will also contain a `console.log` containing the formatted log output.

`timestamp: <message_type-verbosity> [system][component_name][component_type]: message`
 - `timestamp` - Time in milliseconds since program start
 - `message_type` - Either `status`, `info`, `warn`, `error`, or `fatal`
 - `verbosity` - Integer. The higher the number the more verbose. Used for filtering by importance
 - `system` - System generating the message. e.g. "Drive", "Lifter", etc.
 - `component_name` - Name of component within system. e.g. "Right Front Motor", "Piston", etc.
 - `component_type` - Type of component e.g. "Talon", "Solenoid", etc.

#### Binaries
The `binaries/` directory contains all versions of the robot code binary, as well as `ledger.json` which contains info about each. Attributes will be things such as git-commit, date built, user-notes. 
