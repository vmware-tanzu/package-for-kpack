#! Copyright 2022 VMware, Inc.
#! SPDX-License-Identifier: Apache-2.0

load("@ytt:data", "data")
load("@ytt:assert", "assert")

data.values.kp_default_repository or assert.fail("missing kp_default_repository")

if data.values.kp_default_repository_secret and (data.values.kp_default_repository_username or data.values.kp_default_repository_password):
  assert.fail("kp_default_repository_secret cannot be used with kp_default_repository_username or kp_default_repository_password")

if !data.values.kp_default_repository_secret and (!data.values.kp_default_repository_username or !data.values.kp_default_repository_password):
  assert.fail("missing kp_default_repository_username and/or kp_default_repository_password")
