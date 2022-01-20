# kpack

[kpack](https://github.com/pivotal/kpack) utilizes unprivileged Kubernetes primitives to provide builds of OCI images as a platform implementation of [Cloud Native Buildpacks](https://buildpacks.io) (CNB).

## Components

* kpack

## Supported Providers

The following tables shows the providers this package can work with.

 | AKS | EKS | vSphere | Docker |
|-----|-----|---------|--------|
| ✅   | ✅   | ✅       | ✅      |

## Configuration

The following configuration values can be set to customize the kpack installation.

### kpack Configuration

| Value                            | Required/Optional | Description                                                                                                                                                                                       |
|----------------------------------|-------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `kp_default_repository`          | Optional          | Docker repository used for builder images and dependencies (i.e. `gcr.io/my-project/my-repo` or `mydockerhubusername/my-repo`. Used with the [kp cli](https://github.com/vmware-tanzu/kpack-cli). |
| `kp_default_repository_username` | Optional          | Username for `kp_default_repository` (Note: use `_json_key` for GCR)                                                                                                                              |
| `kp_default_repository_password` | Optional          | Password for `kp_default_repository`                                                                                                                                                              |
| `http_proxy`                     | Optional          | The HTTP proxy to use for network traffic                                                                                                                                                         |
| `https_proxy`                    | Optional          | The HTTPS proxy to use for network traffic                                                                                                                                                        |
| `no_proxy`                       | Optional          | A comma-separated list of hostnames, IP addresses, or IP ranges in CIDR format that should not use a proxy                                                                                        |

## Usage Example

### Getting started with kpack using the `kp` cli

> Note: This guide assumes that you have provided the `kp_default_repository` values during install.
> If you do no want to set those values, you can follow the [kubectl flow](#getting-started-with-kpack-using-kubectl) below.
> If you forgot to provide them at install, you can update your values and then update the install using the `tanzu` cli.

> Note: This guide assumes that [`kp`](https://github.com/vmware-tanzu/kpack-cli) has been downloaded. It can be installed using [homebrew](https://github.com/vmware-tanzu/homebrew-kpack-cli), by downloading from the [github release](https://github.com/vmware-tanzu/kpack-cli/releases), or through a [docker image](https://hub.docker.com/r/kpack/kp).

1. Log in to the `kp_default_repository` locally.

    ```bash
    docker login <REGISTRY-HOSTNAME> -u <REGISTRY-USER>
    ```

   > Note: The `<REGISTRY-HOSTNAME>` must be the registry prefix for its corresponding registry
   > - For [dockerhub](https://hub.docker.com/) this should be `https://index.docker.io/v1/`. `kp` also offers a simplified way to create a dockerhub secret with a `--dockerhub` flag.
   > - For [GCR](https://cloud.google.com/container-registry/) this should be `gcr.io`. If you use GCR then the username can be `_json_key` and the password can be the JSON credentials you get from the GCP UI (under `IAM -> Service Accounts` create an account or edit an existing one and create a key with type JSON). `kp` also offers a simplified way to create a gcr secret with a `--gcr` flag.

2. Create a cluster store

   A store resource is a repository of [buildpacks](http://buildpacks.io/) packaged
   in [buildpackages](https://buildpacks.io/docs/buildpack-author-guide/package-a-buildpack/) that can be used by kpack
   to build OCI images. Later in this tutorial, you will reference this store in a Builder configuration.

   We recommend starting with buildpacks from the [paketo project](https://github.com/paketo-buildpacks). The example
   below pulls in java buildpack from the paketo project. If you would like to use a different language you can select a different buildpack.

    ```
    kp clusterstore save default -b gcr.io/paketo-buildpacks/java
    ```

   > Note: Buildpacks are packaged and distributed as buildpackages which are OCI images available on a docker registry. Buildpackages for other languages are available from [paketo](https://github.com/paketo-buildpacks).

3. Create a cluster stack

   A stack resource is the specification for
   a [cloud native buildpacks stack](https://buildpacks.io/docs/concepts/components/stack/) used during build and in the
   resulting app image.

   We recommend starting with the [paketo base stack](https://github.com/paketo-buildpacks/stacks) as shown below:

    ```
    kp clusterstack save base --build-image paketobuildpacks/build:base-cnb --run-image paketobuildpacks/run:base-cnb
    ```

4. Create a Builder

   A Builder is the kpack configuration for a [builder image](https://buildpacks.io/docs/concepts/components/builder/)
   that includes the stack and buildpacks needed to build an OCI image from your app source code.

   The Builder configuration will write to the registry with the secret configured in step one and will reference the
   stack and store created in step three and four. The builder order will determine the order in which buildpacks are
   used in the builder.

    ```
    kp builder save my-builder \
      --tag <IMAGE-TAG> \
      --stack base \
      --store default \
      --buildpack paketo-buildpacks/java \
      -n default
    ```

    - Replace `<IMAGE-TAG>` with a valid, writeable image tag that exists in the same registry as the `kp_default_repository`. The tag should be something like: `your-name/builder` or `gcr.io/your-project/builder`

5. Create a secret with push credentials for the docker registry that you plan on publishing OCI images to with kpack.

   The easiest way to do that is with `kp secret save`

    ```bash
    kp secret save tutorial-registry-credentials \
       --registry <REGISTRY-HOSTNAME> \
       --registry-user <REGISTRY-USER> \
       -n default
    ```

   > Note: The `<REGISTRY-HOSTNAME>` must be the registry prefix for its corresponding registry
   > - For [dockerhub](https://hub.docker.com/) this should be `https://index.docker.io/v1/`. `kp` also offers a simplified way to create a dockerhub secret with a `--dockerhub` flag.
   > - For [GCR](https://cloud.google.com/container-registry/) this should be `gcr.io`. If you use GCR then the username can be `_json_key` and the password can be the JSON credentials you get from the GCP UI (under `IAM -> Service Accounts` create an account or edit an existing one and create a key with type JSON). `kp` also offers a simplified way to create a gcr secret with a `--gcr` flag.

   Your secret create should look something like this:

    ```bash
    kp secret save tutorial-registry-credentials \
       --registry https://index.docker.io/v1/ \
       --registry-user my-dockerhub-username \
       -n default
    ```

   > Note: Learn more about kpack secrets with the [kpack secret documentation](https://github.com/pivotal/kpack/blob/main/docs/secrets.md)


6. Create a kpack image resource

   An image resource is the specification for an OCI image that kpack should build and manage.

   We will create a sample image resource that builds with the builder created in step #4.

   The example included here utilizes
   the [Spring Pet Clinic sample app](https://github.com/spring-projects/spring-petclinic). We encourage you to
   substitute it with your own application.

   Create an image resource:

    ```yaml
    kp image save tutorial-image \
      --tag <IMAGE-TAG> \
      --git https://github.com/spring-projects/spring-petclinic \
      --git-revision 82cb521d636b282340378d80a6307a08e3d4a4c4 \
      --builder my-builder \
      -n default
    ```

    - Make sure to replace `<IMAGE-TAG>` with the tag in the registry of the secret you configured in step #5. Something like:
      `your-name/app` or `gcr.io/your-project/app`
    - If you are using your application source, replace `--git` & `--git-revision`.
   > Note: To use a private git repo follow the instructions in [secrets](https://github.com/pivotal/kpack/blob/main/docs/secrets.md)

   You can now check the status of the image resource.

   ```bash
   kp image status tutorial-image -n default
   ```

   You should see that the image resource has a status Building as it is currently building.

    ```
    Status:         Building
    Message:        --
    LatestImage:    --
    
    Source
    Type:        GitUrl
    Url:         https://github.com/spring-projects/spring-petclinic
    Revision:    82cb521d636b282340378d80a6307a08e3d4a4c4
    
    Builder Ref
    Name:    base
    Kind:    Builder
    
    Last Successful Build
    Id:              --
    Build Reason:    --
    
    Last Failed Build
    Id:              --
    Build Reason:    --
    ```

   You can tail the logs for image resource that is currently building using
   the [kp cli](https://github.com/vmware-tanzu/kpack-cli/blob/main/docs/kp_build_logs.md)

    ```
    kp build logs tutorial-image -n default
    ``` 

   Once the image resource finishes building you can get the fully resolved built OCI image with `kp`

    ```bash
    kp image status tutorial-image -n default
    ```

   The output should look something like this:
    ```
    Status:         Ready
    Message:        --
    LatestImage:    index.docker.io/your-project/app@sha256:6744b...
    
    Source
    Type:        GitUrl
    Url:         https://github.com/spring-projects/spring-petclinic
    Revision:    82cb521d636b282340378d80a6307a08e3d4a4c4
    
    Builder Ref
    Name:    base
    Kind:    Builder
    
    Last Successful Build
    Id:              1
    Build Reason:    BUILDPACK
    Git Revision:    82cb521d636b282340378d80a6307a08e3d4a4c4
    
    BUILDPACK ID                           BUILDPACK VERSION    HOMEPAGE
    paketo-buildpacks/ca-certificates      2.4.0                https://github.com/paketo-buildpacks/ca-certificates
    paketo-buildpacks/bellsoft-liberica    8.4.0                https://github.com/paketo-buildpacks/bellsoft-liberica
    paketo-buildpacks/gradle               5.5.0                https://github.com/paketo-buildpacks/gradle
    paketo-buildpacks/executable-jar       5.2.0                https://github.com/paketo-buildpacks/executable-jar
    paketo-buildpacks/apache-tomcat        6.1.0                https://github.com/paketo-buildpacks/apache-tomcat
    paketo-buildpacks/dist-zip             4.2.0                https://github.com/paketo-buildpacks/dist-zip
    paketo-buildpacks/spring-boot          4.5.0                https://github.com/paketo-buildpacks/spring-boot
    
    Last Failed Build
    Id:              --
    Build Reason:    --
    ```

   The latest built OCI image is available to be used locally via `docker pull` and in a Kubernetes deployment.

7. Run the built app locally

   Download the latest built OCI image available in step #6 and run it with docker.

   ```bash
   docker run -p 8080:8080 <latest-image-with-digest>
   ```

   You should see the java app start up:
   ```
       
              |\      _,,,--,,_
             /,`.-'`'   ._  \-;;,_
    _______ __|,4-  ) )_   .;.(__`'-'__     ___ __    _ ___ _______
    |       | '---''(_/._)-'(_\_)   |   |   |   |  |  | |   |       |
    |    _  |    ___|_     _|       |   |   |   |   |_| |   |       | __ _ _
    |   |_| |   |___  |   | |       |   |   |   |       |   |       | \ \ \ \
    |    ___|    ___| |   | |      _|   |___|   |  _    |   |      _|  \ \ \ \
    |   |   |   |___  |   | |     |_|       |   | | |   |   |     |_    ) ) ) )
    |___|   |_______| |___| |_______|_______|___|_|  |__|___|_______|  / / / /
    ==================================================================/_/_/_/
    
    :: Built with Spring Boot :: 2.2.2.RELEASE
   ``` 

8. kpack rebuilds

   We recommend updating the kpack image resource with a CI/CD tool when new commits are ready to be built.
   > Note: You can also provide a branch or tag as the `spec.git.revision` and kpack will poll and rebuild on updates!

   We can simulate an update from a CI/CD tool by updating the `spec.git.revision` on the image resource configured in step #6.

   If you are using your own application please push an updated commit and use the new commit sha. If you are using
   Spring Pet Clinic you can update the revision to: `4e1f87407d80cdb4a5a293de89d62034fdcbb847`.

   Edit the image resource with:
   ```
   kp image save tutorial-image --git-revision 4e1f87407d80cdb4a5a293de89d62034fdcbb847 -n default
   ``` 

   You should see kpack schedule a new build by running:
   ```
   kp build list tutorial-image -n default
   ``` 
   You should see a new build with

   ```
   BUILD    STATUS     IMAGE                                            REASON
   1        SUCCESS    index.docker.io/your-name/app@sha256:6744b...    BUILDPACK
   2        BUILDING                                                    CONFIG
   ```

   You can tail the logs for the image resource with the kp cli used in step #6.

   ```
   kp build logs tutorial-image -n default
   ```

   > Note: This second build should be notably faster because the buildpacks can leverage the cache from the previous build.

9. Next steps

   The next time new buildpacks are added to the store, kpack will automatically rebuild the builder. If the updated
   buildpacks were used by the tutorial image resource, kpack will automatically create a new build to rebuild your OCI image.


### Getting started with kpack using `kubectl`

1. Create a secret with push credentials for the docker registry that you plan on publishing OCI images to with kpack.

   The easiest way to do that is with `kubectl secret create docker-registry`

    ```bash
    kubectl create secret docker-registry tutorial-registry-credentials \
        --docker-username=user \
        --docker-password=password \
        --docker-server=string \
        --namespace default
    ```

   > Note: The docker server must be the registry prefix for its corresponding registry. For [dockerhub](https://hub.docker.com/) this should be `https://index.docker.io/v1/`.
   For [GCR](https://cloud.google.com/container-registry/) this should be `gcr.io`. If you use GCR then the username can be `_json_key` and the password can be the JSON credentials you get from the GCP UI (under `IAM -> Service Accounts` create an account or edit an existing one and create a key with type JSON).

   Your secret create should look something like this:

    ```bash
    kubectl create secret docker-registry tutorial-registry-credentials \
        --docker-username=my-dockerhub-username \
        --docker-password=my-dockerhub-password \
        --docker-server=https://index.docker.io/v1/ \
        --namespace default
    ```

   > Note: Learn more about kpack secrets with the [kpack secret documentation](secrets.md)

2. Create a service account that references the registry secret created above

    ```yaml
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: tutorial-service-account
      namespace: default
    secrets:
    - name: tutorial-registry-credentials
    imagePullSecrets:
    - name: tutorial-registry-credentials
    ```

   Apply that service account to the cluster

     ```bash
     kubectl apply -f service-account.yaml
     ```

3. Create a cluster store configuration

   A store resource is a repository of [buildpacks](http://buildpacks.io/) packaged in [buildpackages](https://buildpacks.io/docs/buildpack-author-guide/package-a-buildpack/) that can be used by kpack to build OCI images. Later in this tutorial, you will reference this store in a Builder configuration.

   We recommend starting with buildpacks from the [paketo project](https://github.com/paketo-buildpacks). The example below pulls in java and nodejs buildpacks from the paketo project.

    ```yaml
    apiVersion: kpack.io/v1alpha2
    kind: ClusterStore
    metadata:
      name: default
    spec:
      sources:
      - image: gcr.io/paketo-buildpacks/java
      - image: gcr.io/paketo-buildpacks/nodejs
    ```

   Apply this store to the cluster

    ```bash
    kubectl apply -f store.yaml
    ```

   > Note: Buildpacks are packaged and distributed as buildpackages which are docker images available on a docker registry. Buildpackages for other languages are available from [paketo](https://github.com/paketo-buildpacks).

4. Create a cluster stack configuration

   A stack resource is the specification for a [cloud native buildpacks stack](https://buildpacks.io/docs/concepts/components/stack/) used during build and in the resulting app image.

   We recommend starting with the [paketo base stack](https://github.com/paketo-buildpacks/stacks) as shown below:

    ```yaml
    apiVersion: kpack.io/v1alpha2
    kind: ClusterStack
    metadata:
      name: base
    spec:
      id: "io.buildpacks.stacks.bionic"
      buildImage:
        image: "paketobuildpacks/build:base-cnb"
      runImage:
        image: "paketobuildpacks/run:base-cnb"
    ```

   Apply this stack to the cluster

    ```bash
    kubectl apply -f stack.yaml
    ```

5. Create a Builder configuration

   A Builder is the kpack configuration for a [builder image](https://buildpacks.io/docs/concepts/components/builder/) that includes the stack and buildpacks needed to build an OCI image from your app source code.

   The Builder configuration will write to the registry with the secret configured in step one and will reference the stack and store created in step three and four. The builder order will determine the order in which buildpacks are used in the builder.

    ```yaml
    apiVersion: kpack.io/v1alpha2
    kind: Builder
    metadata:
      name: my-builder
      namespace: default
    spec:
      serviceAccountName: tutorial-service-account
      tag: <DOCKER-IMAGE-TAG>
      stack:
        name: base
        kind: ClusterStack
      store:
        name: default
        kind: ClusterStore
      order:
      - group:
        - id: paketo-buildpacks/java
      - group:
        - id: paketo-buildpacks/nodejs
    ```

    - Replace `<DOCKER-IMAGE-TAG>` with a valid image tag that exists in the registry you configured with the `--docker-server` flag when creating a Secret in step #1. The tag should be something like: `your-name/builder` or `gcr.io/your-project/builder`

   Apply this builder to the cluster

     ```bash
     kubectl apply -f builder.yaml
     ```

6. Apply a kpack image resource

   An image resource is the specification for an OCI image that kpack should build and manage.

   We will create a sample image resource that builds with the builder created in step five.

   The example included here utilizes the [Spring Pet Clinic sample app](https://github.com/spring-projects/spring-petclinic). We encourage you to substitute it with your own application.

   Create an image resource:

    ```yaml
    apiVersion: kpack.io/v1alpha2
    kind: Image
    metadata:
      name: tutorial-image
      namespace: default
    spec:
      tag: <DOCKER-IMAGE-TAG>
      serviceAccountName: tutorial-service-account
      builder:
        name: my-builder
        kind: Builder
      source:
        git:
          url: https://github.com/spring-projects/spring-petclinic
          revision: 82cb521d636b282340378d80a6307a08e3d4a4c4
    ```

    - Replace `<DOCKER-IMAGE-TAG>` with a valid image tag that exists in the registry you configured with the `--docker-server` flag when creating a Secret in step #1. Something like: `your-name/app` or `gcr.io/your-project/app`
    - If you are using your application source, replace `source.git.url` & `source.git.revision`.
   > Note: To use a private git repo follow the instructions in [secrets](secrets.md)

   Apply that image resource to the cluster

    ```bash
    kubectl apply -f image.yaml
    ```

   You can now check the status of the image resource.

   ```bash
   kubectl -n default get images
   ```

   You should see that the image resource has an unknown READY status as it is currently building.

   ```
    NAME                  LATESTIMAGE   READY
    tutorial-image                      Unknown
    ```

   You can tail the logs for the image that is currently building using the [kp cli](https://github.com/vmware-tanzu/kpack-cli/blob/main/docs/kp_build_logs.md)

    ```
    kp build logs tutorial-image -n default
    ``` 

   Once the image resource finishes building you can get the fully resolved built OCI image with `kubectl get`

    ```
    kubectl -n default get image tutorial-image
    ```  

   The output should look something like this:
    ```
    NAMESPACE   NAME                  LATESTIMAGE                                        READY
    default     tutorial-image        index.docker.io/your-project/app@sha256:6744b...   True
    ```

   The latest OCI image is available to be used locally via `docker pull` and in a Kubernetes deployment.

8. Run the built app locally

   Download the latest OCI image available in step #6 and run it with docker.

   ```bash
   docker run -p 8080:8080 <latest-image-with-digest>
   ```

   You should see the java app start up:
   ```
       
              |\      _,,,--,,_
             /,`.-'`'   ._  \-;;,_
    _______ __|,4-  ) )_   .;.(__`'-'__     ___ __    _ ___ _______
    |       | '---''(_/._)-'(_\_)   |   |   |   |  |  | |   |       |
    |    _  |    ___|_     _|       |   |   |   |   |_| |   |       | __ _ _
    |   |_| |   |___  |   | |       |   |   |   |       |   |       | \ \ \ \
    |    ___|    ___| |   | |      _|   |___|   |  _    |   |      _|  \ \ \ \
    |   |   |   |___  |   | |     |_|       |   | | |   |   |     |_    ) ) ) )
    |___|   |_______| |___| |_______|_______|___|_|  |__|___|_______|  / / / /
    ==================================================================/_/_/_/
    
    :: Built with Spring Boot :: 2.2.2.RELEASE
   ``` 

9. Rebuilding kpack Images

   We recommend updating the kpack image resource with a CI/CD tool when new commits are ready to be built.
   > Note: You can also provide a branch or tag as the `spec.git.revision` and kpack will poll and rebuild on updates!

   We can simulate an update from a CI/CD tool by updating the `spec.git.revision` on the image resource used in step #6.

   If you are using your own application push an updated commit and use the new commit sha. If you are using Spring Pet Clinic you can update the revision to: `4e1f87407d80cdb4a5a293de89d62034fdcbb847`.

   Edit the image resource with:
   ```
   kubectl -n default edit image tutorial-image 
   ``` 

   You should see kpack schedule a new build by running:
   ```
   kubectl -n default get builds
   ``` 
   You should see a new build with

   ```
   NAME                                IMAGE                                          SUCCEEDED
   tutorial-image-build-1-8mqkc       index.docker.io/your-name/app@sha256:6744b...   True
   tutorial-image-build-2-xsf2l                                                       Unknown
   ```

   You can tail the logs for the image with the kp cli used in step #6.

   ```
   kp build logs tutorial-image -n default -b 2
   ```

   > Note: This second build should be notably faster because the buildpacks can leverage the cache from the previous build.

10. Next steps

    The next time new buildpacks are added to the store, kpack will automatically rebuild the builder. If the updated buildpacks were used by the tutorial image resource, kpack will automatically create a new build to rebuild your OCI image.
