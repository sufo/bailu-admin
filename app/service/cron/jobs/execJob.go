/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package jobs

import "context"

// 新添加的job 必须按照以下格式定义，并实现Exec函数
type ExamplesOne struct {
}

func (t *ExamplesOne) Invoke(ctx context.Context, args map[string]any) (result string, err error) {
	//TODO

	return "", nil
}

func (t *ExamplesOne) Name() string {
	return "测试函数"
}
