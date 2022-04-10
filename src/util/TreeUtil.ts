import { NoteResponsePayload, TreeNode } from "../reducer/noteSlice";

class TreeUtil {
  static parse(seedTree: TreeNode, notes: NoteResponsePayload[], cache?: boolean): TreeNode {
    const tree: TreeNode = notes.reduce((r, n) => {
      const pathArray = n.path.split('/');
      const fileName = pathArray.pop() || "";
      const final = pathArray.reduce((o, name) => {
        let temp = (o.children = o.children || []).find(q => q.name === name);
        if (!temp) o.children.push(temp = {
          name,
          path: o.path ? o.path + '/' + name : name,
          is_dir: true,
          cached: !!cache
        });
        cache != null && (temp.cached = cache);
        o.children.sort((a, b) => (Number(b.is_dir) - Number(a.is_dir)) || a.path.localeCompare(b.path))
        return temp;
      }, r);

      const file = { ...n, name: fileName, cached: !!n.content }
      final.children = final.children || [];
      const index = final.children.findIndex(o => o.path === n.path);
      index > -1 && (final.children[index] = file) || final.children.push(file);
      final.children.sort((a, b) => (Number(b.is_dir) - Number(a.is_dir)) || a.path.localeCompare(b.path))
      cache != null && (final.cached = cache)
      return r;
    }, { ...seedTree });

    return tree;
  }

  static searchNode(root: TreeNode, path: string): TreeNode | null {
    if (root.path == path) {
      return root;
    }

    if (root.children != null) {
      let result = null;
      for (let i = 0; result == null && i < root.children.length; i++) {
        result = TreeUtil.searchNode(root.children[i], path);
      }
      return result;
    }
    return null;
  }

  static deleteNode(root: TreeNode, path: string) {
    if (!root.children) {
      return;
    }
    for (let i = 0; i < root.children.length; i++) {
      const child = root.children[i];
      if (child.path == path) {
        root.children.splice(i, 1);
        break;
      }
      TreeUtil.deleteNode(child, path);
      if (child.is_dir && child.children?.length === 0) {
        // remove empty parent directories on delete
        root.children.splice(i, 1);
      }
    }
  }

  static getChildDirs(tree: TreeNode, path: string): string[] {
    const node = TreeUtil.searchNode(tree, path)
    if (!node?.children) {
      return [];
    }
    return node.children.filter(c => c.is_dir).map(c => c.name);
  }
}

export default TreeUtil;
