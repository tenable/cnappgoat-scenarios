# CNAPPgoat Scenarios Repository

<div align="center">
<img src="https://github.com/ermetic-research/cnappgoat/blob/main/images/logo.png?raw=true" width="40%" alt="CNAPPGoat logo: a smiling goat with a purple cloud background">
</div>

This repository provides a comprehensive collection of Pulumi scenarios utilized
by [CNAPPgoat](https://github.com/ermetic-research/cnappgoat), a multi-cloud, vulnerable-by-design environment
deployment tool – specifically engineered to facilitate practice arenas for defenders and pentesters.

## Scenario Structure

Each scenario in this repository is structured with a `Pulumi.yaml` file and a `main.go` file,
`go.mod` and `go.sum` files should also be included in the scenario.

### The Project File: `Pulumi.yaml

The `Pulumi.yaml` file is a project file in Pulumi and it specifies the runtime, the name of the project,
its description and all parameters used by CNAPPgoat.

Here is an example of a `Pulumi.yaml` file:

```yaml
name: ciem-aws-iam-external-id-3rd-party-role
runtime: go
description: The scenario creates an IAM role without an external ID parameter.
  This exposes your account to confused deputy attacks.
  To fix, include the condition "sts:ExternalId" in your
  IAM role trust policy during the creation process.
cnappgoat-params:
  description: The scenario involves the creation of an IAM role for a 3rd party without
    an external ID parameter. This flaw escalates the risk of impersonation attacks,
    especially 'confused deputy' scenarios, where rogue actors might access your account
    using the same vendor's services. To mitigate this issue, it's essential to include
    an 'External ID' in the trust policy of the IAM role, thereby adding an extra
    layer of security against potential impersonators.
  friendlyName: IAM Role Without External ID
  id: ciem-aws-iam-external-id-3rd-party-role
  module: ciem
  scenarioType: native
  platform: aws
```

Let's break down this example to learn about the different fields:

* `name`: The name of the scenario. This name is used by Pulumi to identify the scenario. 
Our naming convention is `<module>-<platform>-<scenario-name-kebab-case>`
* `runtime`: The programming language used by the scenario, you can use whatever [Pulumi
supported runtime](https://www.pulumi.com/docs/concepts/projects/project-file/#runtime-options) you prefer: 
`nodejs`, `python`, `go`, `dotnet`, `java` or `yaml`
* `description`: A description of the scenario. This description is used by Pulumi to describe the scenario. **Note:**
  **This description can only be up top 256 characters long!**
* `cnappgoat-params`: these are the parameters that CNAPPgoat uses.
  * `description`: A description of the scenario. This description is used by CNAPPgoat to describe the scenario.
Its length is **not** limited
  * `friendlyName`: A friendly name of the scenario. This name is used by CNAPPgoat to describe the scenario. The
convention for this name is `Capitalized Scenario Name`
  * `id`: The ID of the scenario. This ID is used by CNAPPgoat to identify the scenario. Our naming convention is
`<module>-<platform>-<scenario-name-kebab-case>`
  * `module`: The module of the scenario. This module is used by CNAPPgoat to identify the module of the scenario.
(e.g. `ciem`, `cspm`, `cwpp`, etc.)
  * `scenarioType`: The type of the scenario. All scenarios in this repository are `native` scenarios.

### The Pulumi Program
In the same directory as the `Pulumi.yaml` file, you should put the Pulumi program. The Pulumi program is a program
written in the programming language specified in the `Pulumi.yaml` file. This program is used by Pulumi to deploy
the scenario.
#### Go
As of now, all scenarios in this repository are written in Go. 
The `main.go` file is the main program for the scenario. It uses the Pulumi SDK to define resources and their
configurations.

In addition, when adding Go scenarios, you should also add a `go.mod` file and a `go.sum` file.

#### Other Best Practices
These program files are just standard Pulumi programs, so you can use the 
[Pulumi documentation](https://www.pulumi.com/docs/) when writing them. There are a few conventions that you should 
follow when writing scenarios for CNAPPgoat:
* **Resource names** should start with `CNAPPGoat` and be as descriptive as possible.
* **Tags**: Whenever possible, tag resources with `{"Cnappgoat": "true"}`, do note that
tags can be case-sensitive.
* **Output values**: Whenever possible, output values that are relevant to the scenario. For now, these are not used
by CNAPPgoat, but it is best practice to reflect them.

Other than that, try to stick to the best practices and conventions of the programming language you are using.

### Testing your scenario
To test your scenario, just put it into your local directory under the proper directory structure,
e.g. `$HOME/cnappgoat/scenarios/<module>/<platform>/<scenario-name>` and run `cnappgoat provision 
--debug <module>-<platform>-<scenario-name>` to deploy it to a sandbox environment you own.

### Containers and Images
To ensure security, CNAPPgoat scenarios in this repository will only used trusted images/containers or images/containers stored
by the Ermetic Research team. If your scenario uses a custom container,
please [contact the project team](mailto:research+cnappgoat@ermetic.com).
## Contact

To email the project team, contact [research+cnappgoat@ermetic.com](mailto:research+cnappgoat@ermetic.com)

## Disclaimer

* CNAPPGoat scenarios are provided "as is" and without any warranty or support.
* This is a beta version. The project is still in development, so we apologize for any growing pains. We will release a
  stable version in the coming weeks.
