# kvm-vm-clone

This is a command-line tool written in Go for cloning KVM (Kernel-based Virtual Machine) virtual machines. It uses `virt-clone` to create a copy of an existing VM and `virt-customize` to set the hostname of the newly cloned VM.

## Features

- Clone an existing KVM virtual machine
- Automatically set the hostname of the cloned VM
- Check if the source VM is running before attempting to clone

## Prerequisites

Before you can use this tool, ensure you have the following installed on your system:

- KVM and related tools (`virt-clone`, `virt-customize`, `virsh`)
- sudo privileges (the tool uses sudo to run KVM commands)

## Installation

Build the tool:

```
$ go build
```

## Usage

To use the tool, run it with the following command-line arguments:

```
$ ./kvm-vm-clone --source source_vm_name --dest new_vm_name
```

Example:

```
$ ./kvm-vm-clone --source ubuntu-vm --dest ubuntu-vm-clone
```

## Process

The tool performs the following steps:

1. Checks if all required commands are available (`sudo`, `virt-clone`, `virt-customize`, `virsh`)
2. Verifies that the source VM is not currently running
3. Clones the source VM using `virt-clone`
4. Sets the hostname of the new VM using `virt-customize`

## Notes

- The tool requires sudo privileges to run KVM commands. Make sure you have the necessary permissions.
- The source VM must be stopped before cloning. The tool will check this and fail if the VM is running.
- The cloned VM will have its hostname set to the name specified by the `--dest` argument.
- If a VM with the destination name already exists, the cloning process will fail.
- The tool uses `virsh list` to check the VM state, which makes it locale-independent and more robust.

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
