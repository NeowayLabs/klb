// Code generated by "stringer -type=NodeType"; DO NOT EDIT

package ast

import "fmt"

const _NodeType_name = "NodeSetenvNodeBlockNodeAssignmentNodeExecAssignNodeImportexecBeginNodeCommandNodePipeNodeRedirectNodeFnInvexecEndexpressionBeginNodeStringExprNodeIntExprNodeVarExprNodeListExprNodeIndexExprNodeConcatExprexpressionEndNodeStringNodeRforkNodeRforkFlagsNodeIfNodeCommentNodeFnDeclNodeReturnNodeBindFnNodeDumpNodeFor"

var _NodeType_index = [...]uint16{0, 10, 19, 33, 47, 57, 66, 77, 85, 97, 106, 113, 128, 142, 153, 164, 176, 189, 203, 216, 226, 235, 249, 255, 266, 276, 286, 296, 304, 311}

func (i NodeType) String() string {
	i -= 1
	if i < 0 || i >= NodeType(len(_NodeType_index)-1) {
		return fmt.Sprintf("NodeType(%d)", i+1)
	}
	return _NodeType_name[_NodeType_index[i]:_NodeType_index[i+1]]
}