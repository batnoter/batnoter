import ArticleOutlinedIcon from '@mui/icons-material/ArticleOutlined';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import FolderOpenOutlinedIcon from '@mui/icons-material/FolderOpenOutlined';
import FolderOutlinedIcon from '@mui/icons-material/FolderOutlined';
import { TreeItem, TreeView } from '@mui/lab';
import { styled, Toolbar } from '@mui/material';
import MuiDrawer from '@mui/material/Drawer';
import React, { ReactElement } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAppSelector } from '../app/hooks';
import { selectNotesTree, TreeNode } from '../reducer/noteSlice';
import { User } from '../reducer/userSlice';
import { getTitleFromFilename } from '../util/util';

interface Props {
  user: User | null
}

export const DRAWER_WIDTH = 240;
const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(
  () => ({
    '& .MuiDrawer-paper': {
      boxSizing: 'border-box',
      position: 'relative',
      whiteSpace: 'nowrap',
      width: DRAWER_WIDTH,
    },
  }),
)

const AppDrawer: React.FC<Props> = (): ReactElement => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const getAllSubpath = (path: string): string[] => {
    const subpath = path.split('/').map((s, i) => path.split('/').slice(0, i + 1).join('/'));
    subpath.push('/'); // add root path
    return subpath;
  }
  const path = decodeURIComponent(searchParams.get('path') || "%2F");
  const [expanded, setExpanded] = React.useState<string[]>(getAllSubpath(path));
  const tree = useAppSelector(selectNotesTree);

  const handleNodeSelect = (e: React.SyntheticEvent, selectedPath: string) => {
    selectedPath.endsWith('.md') ? navigate("/edit/" + encodeURIComponent(selectedPath))
      : navigate("/?path=" + encodeURIComponent(selectedPath));
  }

  const renderTree = (t: TreeNode) => {
    return (
      <TreeItem key={t.path} nodeId={t.path || "/"} label={getTitleFromFilename(t.name)}
        endIcon={<ArticleOutlinedIcon />} expandIcon={<FolderOutlinedIcon />} collapseIcon={<FolderOpenOutlinedIcon />}>
        {Array.isArray(t.children) ? t.children.map((c) => renderTree(c)) : null}
      </TreeItem>
    )
  }
  const treeJSX = renderTree(tree);

  return (
    <Drawer variant="permanent">
      <Toolbar />
      <TreeView defaultCollapseIcon={<ExpandMoreIcon />} defaultExpandIcon={<ChevronRightIcon />}
        expanded={expanded} selected={path} onNodeSelect={handleNodeSelect}
        onNodeToggle={(e, ids) => setExpanded(ids)} sx={{ flexGrow: 1, minWidth: "max-content", width: "100%" }}>
        {treeJSX}
      </TreeView>
    </Drawer>
  )
}

export default AppDrawer;
