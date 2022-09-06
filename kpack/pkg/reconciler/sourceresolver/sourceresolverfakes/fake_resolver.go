// Code generated by counterfeiter. DO NOT EDIT.
package sourceresolverfakes

import (
	"context"
	"sync"

	"github.com/pivotal/kpack/pkg/apis/build/v1alpha2"
	"github.com/pivotal/kpack/pkg/apis/core/v1alpha1"
	"github.com/pivotal/kpack/pkg/reconciler/sourceresolver"
)

type FakeResolver struct {
	CanResolveStub        func(*v1alpha2.SourceResolver) bool
	canResolveMutex       sync.RWMutex
	canResolveArgsForCall []struct {
		arg1 *v1alpha2.SourceResolver
	}
	canResolveReturns struct {
		result1 bool
	}
	canResolveReturnsOnCall map[int]struct {
		result1 bool
	}
	ResolveStub        func(context.Context, *v1alpha2.SourceResolver) (v1alpha1.ResolvedSourceConfig, error)
	resolveMutex       sync.RWMutex
	resolveArgsForCall []struct {
		arg1 context.Context
		arg2 *v1alpha2.SourceResolver
	}
	resolveReturns struct {
		result1 v1alpha1.ResolvedSourceConfig
		result2 error
	}
	resolveReturnsOnCall map[int]struct {
		result1 v1alpha1.ResolvedSourceConfig
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeResolver) CanResolve(arg1 *v1alpha2.SourceResolver) bool {
	fake.canResolveMutex.Lock()
	ret, specificReturn := fake.canResolveReturnsOnCall[len(fake.canResolveArgsForCall)]
	fake.canResolveArgsForCall = append(fake.canResolveArgsForCall, struct {
		arg1 *v1alpha2.SourceResolver
	}{arg1})
	fake.recordInvocation("CanResolve", []interface{}{arg1})
	fake.canResolveMutex.Unlock()
	if fake.CanResolveStub != nil {
		return fake.CanResolveStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.canResolveReturns
	return fakeReturns.result1
}

func (fake *FakeResolver) CanResolveCallCount() int {
	fake.canResolveMutex.RLock()
	defer fake.canResolveMutex.RUnlock()
	return len(fake.canResolveArgsForCall)
}

func (fake *FakeResolver) CanResolveCalls(stub func(*v1alpha2.SourceResolver) bool) {
	fake.canResolveMutex.Lock()
	defer fake.canResolveMutex.Unlock()
	fake.CanResolveStub = stub
}

func (fake *FakeResolver) CanResolveArgsForCall(i int) *v1alpha2.SourceResolver {
	fake.canResolveMutex.RLock()
	defer fake.canResolveMutex.RUnlock()
	argsForCall := fake.canResolveArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeResolver) CanResolveReturns(result1 bool) {
	fake.canResolveMutex.Lock()
	defer fake.canResolveMutex.Unlock()
	fake.CanResolveStub = nil
	fake.canResolveReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeResolver) CanResolveReturnsOnCall(i int, result1 bool) {
	fake.canResolveMutex.Lock()
	defer fake.canResolveMutex.Unlock()
	fake.CanResolveStub = nil
	if fake.canResolveReturnsOnCall == nil {
		fake.canResolveReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.canResolveReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeResolver) Resolve(arg1 context.Context, arg2 *v1alpha2.SourceResolver) (v1alpha1.ResolvedSourceConfig, error) {
	fake.resolveMutex.Lock()
	ret, specificReturn := fake.resolveReturnsOnCall[len(fake.resolveArgsForCall)]
	fake.resolveArgsForCall = append(fake.resolveArgsForCall, struct {
		arg1 context.Context
		arg2 *v1alpha2.SourceResolver
	}{arg1, arg2})
	fake.recordInvocation("Resolve", []interface{}{arg1, arg2})
	fake.resolveMutex.Unlock()
	if fake.ResolveStub != nil {
		return fake.ResolveStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.resolveReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeResolver) ResolveCallCount() int {
	fake.resolveMutex.RLock()
	defer fake.resolveMutex.RUnlock()
	return len(fake.resolveArgsForCall)
}

func (fake *FakeResolver) ResolveCalls(stub func(context.Context, *v1alpha2.SourceResolver) (v1alpha1.ResolvedSourceConfig, error)) {
	fake.resolveMutex.Lock()
	defer fake.resolveMutex.Unlock()
	fake.ResolveStub = stub
}

func (fake *FakeResolver) ResolveArgsForCall(i int) (context.Context, *v1alpha2.SourceResolver) {
	fake.resolveMutex.RLock()
	defer fake.resolveMutex.RUnlock()
	argsForCall := fake.resolveArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeResolver) ResolveReturns(result1 v1alpha1.ResolvedSourceConfig, result2 error) {
	fake.resolveMutex.Lock()
	defer fake.resolveMutex.Unlock()
	fake.ResolveStub = nil
	fake.resolveReturns = struct {
		result1 v1alpha1.ResolvedSourceConfig
		result2 error
	}{result1, result2}
}

func (fake *FakeResolver) ResolveReturnsOnCall(i int, result1 v1alpha1.ResolvedSourceConfig, result2 error) {
	fake.resolveMutex.Lock()
	defer fake.resolveMutex.Unlock()
	fake.ResolveStub = nil
	if fake.resolveReturnsOnCall == nil {
		fake.resolveReturnsOnCall = make(map[int]struct {
			result1 v1alpha1.ResolvedSourceConfig
			result2 error
		})
	}
	fake.resolveReturnsOnCall[i] = struct {
		result1 v1alpha1.ResolvedSourceConfig
		result2 error
	}{result1, result2}
}

func (fake *FakeResolver) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.canResolveMutex.RLock()
	defer fake.canResolveMutex.RUnlock()
	fake.resolveMutex.RLock()
	defer fake.resolveMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeResolver) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ sourceresolver.Resolver = new(FakeResolver)
