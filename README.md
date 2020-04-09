# Cloud Pak LifeCycle Operator - cpeir project

## Overview

The IBM Cloud Pak lifecycle management operator will run as native extension to Kuberentes (OpenShift)
to provide insight and manage IBM Cloud Pak components. Initially this operator will check whether
a specified IBM Cloud Pak can be installed in the cluster. It is expected that this tool will be extended over time
to accommodate the whole lifecycle management of IBM Cloud Pak, including:

  - manage a Cloud Pak deployment
  - check for requisite resource requirements
  - verify the installation
  - do stability check
  - check for updates

The documentation for this tool contains the following:

- [Technical Design](docs/design.md)

    Older docs:
    - [Technical Design](docs/design/technical-design.md)
    - [Minimum Viable Product (MVP) Iteration 1](docs/design/mvp01.md)
- [Tool development guide](docs/development.md)
- [Installation and usage](docs/installation.md)
- [Usage](docs/usage.md)
