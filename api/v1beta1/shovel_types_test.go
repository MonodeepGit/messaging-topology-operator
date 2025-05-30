package v1beta1

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Shovel spec", func() {
	var (
		namespace = "default"
		ctx       = context.Background()
	)

	It("creates a shovel with minimal configurations", func() {
		shovel := Shovel{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-shovel",
				Namespace: namespace,
			},
			Spec: ShovelSpec{
				Name: "test-shovel",
				RabbitmqClusterReference: RabbitmqClusterReference{
					Name: "some-cluster",
				},
				UriSecret: &corev1.LocalObjectReference{
					Name: "a-secret",
				},
			}}
		Expect(k8sClient.Create(ctx, &shovel)).To(Succeed())
		fetched := &Shovel{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{
			Name:      shovel.Name,
			Namespace: shovel.Namespace,
		}, fetched)).To(Succeed())
		Expect(fetched.Spec.Name).To(Equal("test-shovel"))
		Expect(fetched.Spec.Vhost).To(Equal("/"))
		Expect(fetched.Spec.RabbitmqClusterReference.Name).To(Equal("some-cluster"))
		Expect(fetched.Spec.UriSecret.Name).To(Equal("a-secret"))
	})

	It("creates shovel with configurations", func() {
		shovel := Shovel{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-shovel-configurations",
				Namespace: namespace,
			},
			Spec: ShovelSpec{
				Name:  "test-shovel-configurations",
				Vhost: "test-vhost",
				RabbitmqClusterReference: RabbitmqClusterReference{
					Name: "some-cluster",
				},
				UriSecret: &corev1.LocalObjectReference{
					Name: "a-secret",
				},
				AckMode:                          "no-ack",
				AddForwardHeaders:                true,
				DeleteAfter:                      "never",
				DestinationAddForwardHeaders:     true,
				DestinationAddTimestampHeader:    true,
				DestinationAddress:               "myQueue",
				DestinationApplicationProperties: &runtime.RawExtension{Raw: []byte(`{"key": "a-property"}`)},
				DestinationMessageAnnotations:    &runtime.RawExtension{Raw: []byte(`{"key": "a-property"}`)},
				DestinationExchange:              "an-exchange",
				DestinationExchangeKey:           "a-key",
				DestinationProperties:            &runtime.RawExtension{Raw: []byte(`{"key": "a-property"}`)},
				DestinationProtocol:              "amqp091",
				DestinationPublishProperties:     &runtime.RawExtension{Raw: []byte(`{"delivery_mode": 1}`)},
				DestinationQueue:                 "a-queue",
				DestinationQueueArgs: &runtime.RawExtension{
					Raw: []byte(`{"x-queue-type": "quorum"}`),
				},
				PrefetchCount:       10,
				ReconnectDelay:      10,
				SourceAddress:       "myQueue",
				SourceDeleteAfter:   "never",
				SourceExchange:      "an-exchange",
				SourceExchangeKey:   "a-key",
				SourcePrefetchCount: 10,
				SourceProtocol:      "amqp091",
				SourceQueue:         "a-queue",
				SourceQueueArgs: &runtime.RawExtension{
					Raw: []byte(`{"x-queue-type": "quorum"}`),
				},
				SourceConsumerArgs: &runtime.RawExtension{Raw: []byte(`{"arg": "arg-value"}`)},
			}}
		Expect(k8sClient.Create(ctx, &shovel)).To(Succeed())
		fetched := &Shovel{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{
			Name:      shovel.Name,
			Namespace: shovel.Namespace,
		}, fetched)).To(Succeed())

		Expect(fetched.Spec.Name).To(Equal("test-shovel-configurations"))
		Expect(fetched.Spec.Vhost).To(Equal("test-vhost"))
		Expect(fetched.Spec.RabbitmqClusterReference.Name).To(Equal("some-cluster"))
		Expect(fetched.Spec.UriSecret.Name).To(Equal("a-secret"))
		Expect(fetched.Spec.AckMode).To(Equal("no-ack"))
		Expect(fetched.Spec.AddForwardHeaders).To(BeTrue())
		Expect(fetched.Spec.DeleteAfter).To(Equal("never"))

		Expect(fetched.Spec.DestinationAddTimestampHeader).To(BeTrue())
		Expect(fetched.Spec.DestinationAddForwardHeaders).To(BeTrue())
		Expect(fetched.Spec.DestinationAddress).To(Equal("myQueue"))
		Expect(fetched.Spec.DestinationApplicationProperties.Raw).To(Equal([]byte(`{"key":"a-property"}`)))
		Expect(fetched.Spec.DestinationExchange).To(Equal("an-exchange"))
		Expect(fetched.Spec.DestinationExchangeKey).To(Equal("a-key"))
		Expect(fetched.Spec.DestinationProperties.Raw).To(Equal([]byte(`{"key":"a-property"}`)))
		Expect(fetched.Spec.DestinationMessageAnnotations.Raw).To(Equal([]byte(`{"key":"a-property"}`)))
		Expect(fetched.Spec.DestinationQueue).To(Equal("a-queue"))
		Expect(fetched.Spec.DestinationQueueArgs.Raw).To(Equal([]byte(`{"x-queue-type":"quorum"}`)))
		Expect(fetched.Spec.PrefetchCount).To(Equal(10))
		Expect(fetched.Spec.ReconnectDelay).To(Equal(10))

		Expect(fetched.Spec.SourceAddress).To(Equal("myQueue"))
		Expect(fetched.Spec.SourceDeleteAfter).To(Equal("never"))
		Expect(fetched.Spec.SourceExchange).To(Equal("an-exchange"))
		Expect(fetched.Spec.SourceExchangeKey).To(Equal("a-key"))
		Expect(fetched.Spec.SourcePrefetchCount).To(Equal(10))
		Expect(fetched.Spec.SourceProtocol).To(Equal("amqp091"))
		Expect(fetched.Spec.SourceQueue).To(Equal("a-queue"))
		Expect(fetched.Spec.SourceQueueArgs.Raw).To(Equal([]byte(`{"x-queue-type":"quorum"}`)))
		Expect(fetched.Spec.SourceConsumerArgs.Raw).To(Equal([]byte(`{"arg":"arg-value"}`)))
	})

	When("creating a shovel with an invalid 'AckMode' value", func() {
		It("fails with validation errors", func() {
			shovel := Shovel{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "an-invalid-ackmode",
					Namespace: namespace,
				},
				Spec: ShovelSpec{
					Name: "an-invalid-ackmode",
					RabbitmqClusterReference: RabbitmqClusterReference{
						Name: "some-cluster",
					},
					UriSecret: &corev1.LocalObjectReference{
						Name: "a-secret",
					},
					AckMode: "an-invalid-ackmode",
				}}
			Expect(k8sClient.Create(ctx, &shovel)).To(HaveOccurred())
			Expect(k8sClient.Create(ctx, &shovel)).To(MatchError(`Shovel.rabbitmq.com "an-invalid-ackmode" is invalid: spec.ackMode: Unsupported value: "an-invalid-ackmode": supported values: "on-confirm", "on-publish", "no-ack"`))
		})
	})

	When("creating a shovel with unsupported protocol", func() {
		It("fails with validation errors", func() {
			shovel := Shovel{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "an-invalid-destprotocol",
					Namespace: namespace,
				},
				Spec: ShovelSpec{
					Name: "an-invalid-destprotocol",
					RabbitmqClusterReference: RabbitmqClusterReference{
						Name: "some-cluster",
					},
					UriSecret: &corev1.LocalObjectReference{
						Name: "a-secret",
					},
					SourceProtocol:      "amqp091",
					DestinationProtocol: "stomp",
				}}
			Expect(k8sClient.Create(ctx, &shovel)).To(HaveOccurred())
			Expect(k8sClient.Create(ctx, &shovel)).To(MatchError(`Shovel.rabbitmq.com "an-invalid-destprotocol" is invalid: spec.destProtocol: Unsupported value: "stomp": supported values: "amqp091", "amqp10"`))
		})

		It("fails with validation errors", func() {
			shovel := Shovel{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "an-invalid-srcprotocol",
					Namespace: namespace,
				},
				Spec: ShovelSpec{
					Name: "an-invalid-srcprotocol",
					RabbitmqClusterReference: RabbitmqClusterReference{
						Name: "some-cluster",
					},
					UriSecret: &corev1.LocalObjectReference{
						Name: "a-secret",
					},
					SourceProtocol:      "mqtt",
					DestinationProtocol: "amqp10",
				}}
			Expect(k8sClient.Create(ctx, &shovel)).To(HaveOccurred())
			Expect(k8sClient.Create(ctx, &shovel)).To(MatchError(`Shovel.rabbitmq.com "an-invalid-srcprotocol" is invalid: spec.srcProtocol: Unsupported value: "mqtt": supported values: "amqp091", "amqp10"`))
		})
	})

	It("creates a shovel with non-default DeletionPolicy", func() {
		shovel := Shovel{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "shovel-with-retain-policy",
				Namespace: namespace,
			},
			Spec: ShovelSpec{
				Name:           "shovel-with-retain-policy",
				DeletionPolicy: "retain",
				RabbitmqClusterReference: RabbitmqClusterReference{
					Name: "some-cluster",
				},
				UriSecret: &corev1.LocalObjectReference{
					Name: "a-secret",
				},
			},
		}
		Expect(k8sClient.Create(ctx, &shovel)).To(Succeed())
		fetched := &Shovel{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{
			Name:      shovel.Name,
			Namespace: shovel.Namespace,
		}, fetched)).To(Succeed())

		Expect(fetched.Spec.DeletionPolicy).To(Equal("retain"))
		Expect(fetched.Spec.Name).To(Equal("shovel-with-retain-policy"))
		Expect(fetched.Spec.RabbitmqClusterReference).To(Equal(RabbitmqClusterReference{
			Name: "some-cluster",
		}))
		Expect(fetched.Spec.UriSecret.Name).To(Equal("a-secret"))
	})
})
