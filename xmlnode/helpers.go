package xmlnode

func Update(id string, attributes map[string]string) *XmlNode {
	if attributes == nil {
		attributes = map[string]string{}
	}
	attributes["id"] = id
	return Node("update", attributes)
}
