package config

const Version = "v1.0"
const Usage = `########## config file: RecommendConfFile ##########
Usage:
Args:
  -h/--help
  --kubeconfig  # specify kubeconfig file, default value is ""
  -d            # specify node name,       default value is []
  -u            # specify user name for ssh worker node, default value is "root"
  -n            # specify namespce name,   default value is []
  -p            # specify pod name,        default value is []
  -c            # specify container name,  default value is []

Commands:
  version:
    ek version                                    # show version
  config:
    ek config                                     # show      configuration in config file
    ek config reset                               # reset all configuration in config file
    ek config -d node1,node2                      # change node name in config file
    ek config -d reset                            # reset  node name in config file, default value is: []
    ek config -u username                         # change user name in config file
    ek config -u reset                            # reset  user name in config file, default value is: "root"
    ek config -n ns1,ns2                          # change namespace name in config file
    ek config -n reset                            # reset  namespace name in config file, default value is: []
    ek config -p pod1,pod2                        # change pod name in config file
    ek config -p reset                            # reset  pod name in config file, default value is: []
    ek config -c container1,container2            # change container name in config file
    ek config -c reset                            # reset  container name in config file, default value is: []
    ek config --kubeconfig=/path/to/kbcon.yaml    # change kubeconfig file path in config file
    ek config --kubeconfig=reset                  # reset  kubeconfig file path in config file, default value is: ""
  get:
    ek get node                  <-d node1,node2 -u username>                # show specific master or worker nodes
    ek get node                  <-d reset -u username>                      # show all      master or worker nodes
    ek get node sriov            <-d node1,node2>                            # show specific master or worker nodes sriov resources
    ek get node sriov            <-d reset>                                  # show all      master or worker nodes sriov resources
    ek get node lab|label        <-d node1,node2>                            # show specific master or worker nodes labels
    ek get node lab|label        <-d reset>                                  # show all      master or worker nodes labels
    ek get node ann|anno         <-d node1,node2>                            # show specific master or worker nodes annotations
    ek get node ann|anno         <-d reset>                                  # show all      master or worker nodes annotations
    ek get ns                    <-n ns1,ns2>                                # show specific namespace
    ek get ns                    <-n reset>                                  # show all      namespace
    ek get ns lab|label          <-n ns1,ns2>                                # show specific namespace labels
    ek get ns lab|label          <-n reset>                                  # show all      namespace labels
    ek get ns ann|anno           <-n ns1,ns2>                                # show specific namespace annotations
    ek get ns ann|anno           <-n reset>                                  # show all      namespace annotations
    ek get po|pod|pods           <-n ns1,ns2 -p pod1,pod2>                   # show specific pods
    ek get po|pod|pods           <-n reset -p reset>                         # show all      pods
    ek get po|pod|pods lab|label <-n ns1,ns2 -p pod1,pod2>                   # show specific pods labels
    ek get po|pod|pods lab|label <-n reset -p reset>                         # show all      pods labels
    ek get po|pod|pods ann|anno  <-n ns1,ns2 -p pod1,pod2>                   # show specific pods annotations
    ek get po|pod|pods ann|anno  <-n reset -p reset>                         # show all      pods annotations
    ek get con|cont|container    <-n ns1,ns2 -p pod1,pod2 -c cont1,cont2>    # show specific container
    ek get con|cont|container    <-n reset -p reset -c reset>                # show all      container
    ek get pvc|PVC               [pvc1,pvc2] <-n ns1,ns2>                    # show specific pvc
    ek get pvc|PVC               <-n reset>                                  # show all      pvc
    ek get hostpath|HostPath     <-n ns1,ns2>                                # show specific HostPath volumes
    ek get res|resource          <-n ns1,ns2 -p pod1,pod2 -c cont1,cont2>    # show specific container cpu/memory/sriov resource
    ek get vol|volume|VOL        <-n ns1,ns2 -p pod1,pod2 -c cont1,cont2>    # show specific container volumes mount, and pods volumes
    ek get env|ENV               <-n ns1,ns2 -p pod1,pod2 -c cont1,cont2>    # show specific container env config in chart
    ek get img|image|images      <-n ns1,ns2 -p pod1,pod2>                   # show specific container images
  exec:
    ek exec        <shell command> <-n -p -c>        # exec shell command, e.g: ls /root
  netcmd:
    ek netcmd      <net command> <-d -u -n -p -c>    # exec net command with nsenter in worker node, e.g: ip, ping, route, netstat
`
