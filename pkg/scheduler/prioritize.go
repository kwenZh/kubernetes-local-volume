package scheduler

import (
	"github.com/kubernetes-local-volume/kubernetes-local-volume/pkg/common/logging"
	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

func (lvs *LocalVolumeScheduler) PrioritizeHandler(args schedulerapi.ExtenderArgs) (*schedulerapi.HostPriorityList, error) {
	return lvs.prioritize(*args.Pod, args.Nodes.Items)
}

func (lvs *LocalVolumeScheduler) prioritize(pod v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error) {
	logger := logging.FromContext(lvs.ctx)
	curMaxFreeSizeNode := lvs.getMaxFreeSizeNode(nodes)
	logger.Infof("local volume scheduler prioritize pod(%s) namespace(%s) max free size node(%s)",
		pod.Name, pod.Namespace, curMaxFreeSizeNode)

	var priorityList schedulerapi.HostPriorityList
	priorityList = make([]schedulerapi.HostPriority, len(nodes))
	for i, node := range nodes {
		priorityList[i] = schedulerapi.HostPriority{
			Host: node.Name,
		}
		if curMaxFreeSizeNode == node.Name {
			priorityList[i].Score = 100
		} else {
			priorityList[i].Score = 0
		}
	}
	return &priorityList, nil
}

func (lvs *LocalVolumeScheduler) getMaxFreeSizeNode(nodes []v1.Node) string {
	var curMax uint64
	var curNode string

	for _, node := range nodes {
		freeSize := lvs.getNodeFreeSize(node.Name)
		if freeSize > curMax {
			curMax = freeSize
			curNode = node.Name
		}
	}
	return curNode
}
