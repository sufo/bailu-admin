/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package dto

type ConfigParams struct {
	Name string `json:"path" query:"name,like"` //
	Key  string `json:"key" query:"key"`
	Type string `json:"type" query:"type"`
}
