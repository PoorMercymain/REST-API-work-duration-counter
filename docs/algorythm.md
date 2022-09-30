- `SELECT * FROM work;`
- `GROUP BY TASK_ID`
- find root
  - group parent_id = root_id 
  - repeat for current_node_id as parent_id
  - until children not found for node
  
- find leafs
  - move_to parent
  - repeat
  - until parent_id != nil

   go
    type TreeNode struct{
        children []*TreeNode
    }
  
    type Node interface{
    AddNext() Node
    Moveto() []Node
  }

- find root 
  -[]work - where task_id = equal 
  - find root - work[[0]]: delete from work 
  - find children - delete children from work
  - repeat
  - until len(work)==0

- find leaf
  - []work - where task_id = equal 
  - []leaf - []work where is_leaf == true
  - move to parent
  - map[[parent]]=child -> node []children
  - delete leaf
  - []leaf = []parent
  - repeat
  - until []parent is nil for []leaf