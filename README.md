# Velero Annotations Prune Plugin

This repository is based on the Velero Plugins Example repository: https://github.com/vmware-tanzu/velero-plugin-example/

The Annotations Prune plugin removes specified annotations (specify part TBD) from the pod spec.

## Kinds of Plugin

This is a **Restore Item Action** - performs arbitrary logic on individual items prior to restoring them in the Kubernetes cluster.

For more information, please see the full [plugin documentation](https://velero.io/docs/main/overview-plugins/).

## Deploying the plugins

To deploy your plugin image to an Velero server:

1. Make sure your image is pushed to a registry that is accessible to your cluster's nodes.
2. Run `velero plugin add <registry/image:version>`. Example with a dockerhub image: `velero plugin add patst/velero-plugin-osm-prune:main`.

## Using the plugins

When the plugin is deployed, it is only made available to use. To make the plugin effective, you must modify your configuration:

Backup/Restore actions:

1. Add the plugin to Velero as described in the Deploying the plugins section. (e.g. `velero plugin add getais/velero-plugin-annotations-prune:master`)
2. The plugin will be used for the next `backup/restore`.
