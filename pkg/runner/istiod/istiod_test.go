package istiod

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsIstiod(t *testing.T) {
	type testcase struct {
		pod      v1.Pod
		expected bool
	}

	testcases := []testcase{
		{
			pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":   "istiod",
						"istio": "pilot",
					},
				},
			},
			expected: true,
		},
		{
			pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"k8s": "pilot",
					},
				},
			},
			expected: false,
		},
	}
	for i, tc := range testcases {
		if v := isIstiod(tc.pod); v != tc.expected {
			t.Errorf("[%d] got %t, expected %t", i, v, tc.expected)
		}
	}
}
