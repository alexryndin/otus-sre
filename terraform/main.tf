provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "minikube"
}

resource "kubernetes_namespace" "otus" {
  metadata {
    name = "otus"
  }
}

resource "kubernetes_role" "pod_reader" {
  metadata {
    name      = "pod-reader"
    namespace = "otus"
  }

  rule {
    api_groups = [""]
    resources  = ["pods"]
    verbs      = ["get", "list", "watch"]
  }
}

resource "kubernetes_service_account" "pods_sa" {
  metadata {
    name      = "pods-sa"
    namespace = "otus"
  }
}

resource "kubernetes_role_binding" "pod_reader_binding" {
  metadata {
    name      = "pod-reader-binding"
    namespace = "otus"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "Role"
    name      = "pod-reader"
  }

  subject {
    kind      = "ServiceAccount"
    name      = "pods-sa"
    namespace = "otus"
  }
}
