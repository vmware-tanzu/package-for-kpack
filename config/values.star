# Copyright 2021 VMware, Inc.
# SPDX-License-Identifier: Apache-2.0

load("@ytt:data", "data")
load("@ytt:base64", "base64")
load("@ytt:json", "json")
load("@ytt:assert", "assert")

kp_default_docker_configjson = ""

if data.values.kp_default_repository != "":
    data.values.kp_default_repository_username or assert.fail("missing kp_default_repository_username")
    data.values.kp_default_repository_password or assert.fail("missing kp_default_repository_password")

    # extract the docker registry from the repository string
    kp_default_registry = "https://index.docker.io/v1/"
    parts = data.values.kp_default_repository.split("/", 1)
    if len(parts) == 2:
        if '.' in parts[0] or ':' in parts[0]:
            kp_default_registry = parts[0]
        end
    elif len(parts) == 1:
        assert.fail("kp_default_repository must be a valid writeable repository and must include a '/'")
    end

    kp_default_docker_auth = base64.encode("{}:{}".format(data.values.kp_default_repository_username, data.values.kp_default_repository_password))
    kp_default_docker_creds = {"username": data.values.kp_default_repository_username, "password": data.values.kp_default_repository_password, "auth": kp_default_docker_auth}
    kp_default_docker_configjson = base64.encode(json.encode({"auths": {kp_default_registry: kp_default_docker_creds}}))
end


