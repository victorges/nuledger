package rules

type CommitFuncSlice []CommitFunc

func (c CommitFuncSlice) Commit() {
	for _, commitFunc := range c {
		if commitFunc != nil {
			commitFunc()
		}
	}
}
