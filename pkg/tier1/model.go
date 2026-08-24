package tier1

import "math"


type Model struct{}
func score(input []float64) []float64 {
    var var0 float64
    if input[2] < 1.0 {
        if input[1] < 2.0 {
            if input[1] < 1.0 {
                if input[3] < 1.0 {
                    var0 = -0.043393604
                } else {
                    var0 = -0.27150258
                }
            } else {
                if input[0] < 2.0 {
                    var0 = -0.24934907
                } else {
                    var0 = -0.044210527
                }
            }
        } else {
            if input[0] < 1.0 {
                if input[1] < 4.0 {
                    var0 = -0.31074607
                } else {
                    var0 = -0.39830148
                }
            } else {
                var0 = -0.3993229
            }
        }
    } else {
        if input[3] < 2.0 {
            if input[1] < 4.0 {
                if input[0] < 2.0 {
                    var0 = 0.24279475
                } else {
                    var0 = 0.3647059
                }
            } else {
                var0 = -0.2769231
            }
        } else {
            var0 = -0.34074077
        }
    }
    var var1 float64
    if input[2] < 1.0 {
        if input[1] < 2.0 {
            if input[1] < 1.0 {
                if input[3] < 2.0 {
                    var1 = -0.050633557
                } else {
                    var1 = -0.3426834
                }
            } else {
                if input[3] < 6.0 {
                    var1 = -0.19265744
                } else {
                    var1 = 0.2676275
                }
            }
        } else {
            if input[0] < 1.0 {
                if input[1] < 4.0 {
                    var1 = -0.25523195
                } else {
                    var1 = -0.33281562
                }
            } else {
                var1 = -0.3335664
            }
        }
    } else {
        if input[3] < 2.0 {
            if input[1] < 4.0 {
                if input[1] < 2.0 {
                    var1 = 0.22266786
                } else {
                    var1 = 0.13723922
                }
            } else {
                var1 = -0.24199434
            }
        } else {
            var1 = -0.29028478
        }
    }
    var var2 float64
    if input[2] < 1.0 {
        if input[1] < 2.0 {
            if input[1] < 1.0 {
                if input[3] < 1.0 {
                    var2 = -0.024708686
                } else {
                    var2 = -0.1974273
                }
            } else {
                if input[0] < 2.0 {
                    var2 = -0.17183106
                } else {
                    var2 = 0.0015023648
                }
            }
        } else {
            if input[0] < 1.0 {
                if input[1] < 4.0 {
                    var2 = -0.21710534
                } else {
                    var2 = -0.29484123
                }
            } else {
                var2 = -0.29553226
            }
        }
    } else {
        if input[3] < 2.0 {
            if input[1] < 4.0 {
                if input[0] < 2.0 {
                    var2 = 0.16293389
                } else {
                    var2 = 0.27874935
                }
            } else {
                var2 = -0.21628742
            }
        } else {
            var2 = -0.25709608
        }
    }
    var var3 float64
    if input[1] < 2.0 {
        if input[2] < 1.0 {
            if input[1] < 1.0 {
                if input[3] < 2.0 {
                    var3 = -0.031207083
                } else {
                    var3 = -0.2796455
                }
            } else {
                if input[3] < 6.0 {
                    var3 = -0.1347407
                } else {
                    var3 = 0.24547787
                }
            }
        } else {
            if input[3] < 1.0 {
                if input[0] < 1.0 {
                    var3 = 0.14134602
                } else {
                    var3 = 0.22214353
                }
            } else {
                var3 = -0.31436837
            }
        }
    } else {
        if input[2] < 1.0 {
            if input[0] < 1.0 {
                if input[1] < 4.0 {
                    var3 = -0.18794376
                } else {
                    var3 = -0.27020755
                }
            } else {
                var3 = -0.27092227
            }
        } else {
            if input[0] < 1.0 {
                if input[1] < 4.0 {
                    var3 = 0.13706076
                } else {
                    var3 = -0.17568006
                }
            } else {
                var3 = -0.37258512
            }
        }
    }
    var var4 float64
    if input[1] < 2.0 {
        if input[2] < 1.0 {
            if input[3] < 2.0 {
                if input[1] < 1.0 {
                    var4 = -0.025062231
                } else {
                    var4 = -0.10118157
                }
            } else {
                if input[3] < 6.0 {
                    var4 = -0.28647614
                } else {
                    var4 = 0.1709879
                }
            }
        } else {
            if input[3] < 1.0 {
                if input[0] < 1.0 {
                    var4 = 0.119492136
                } else {
                    var4 = 0.19605541
                }
            } else {
                var4 = -0.27959797
            }
        }
    } else {
        if input[2] < 1.0 {
            if input[0] < 1.0 {
                if input[1] < 4.0 {
                    var4 = -0.16399205
                } else {
                    var4 = -0.25311586
                }
            } else {
                var4 = -0.25390285
            }
        } else {
            if input[0] < 1.0 {
                if input[1] < 3.0 {
                    var4 = 0.019196415
                } else {
                    var4 = 0.17568128
                }
            } else {
                var4 = -0.31431565
            }
        }
    }
    var var5 float64
    if input[1] < 2.0 {
        if input[2] < 1.0 {
            if input[3] < 2.0 {
                if input[1] < 1.0 {
                    var5 = -0.020122295
                } else {
                    var5 = -0.08435035
                }
            } else {
                if input[3] < 6.0 {
                    var5 = -0.26422733
                } else {
                    var5 = 0.14054374
                }
            }
        } else {
            if input[3] < 1.0 {
                if input[0] < 2.0 {
                    var5 = 0.11196852
                } else {
                    var5 = 0.23353724
                }
            } else {
                var5 = -0.25286192
            }
        }
    } else {
        if input[0] < 1.0 {
            if input[2] < 1.0 {
                if input[1] < 4.0 {
                    var5 = -0.14337347
                } else {
                    var5 = -0.24070421
                }
            } else {
                if input[1] < 4.0 {
                    var5 = 0.101465896
                } else {
                    var5 = -0.17568016
                }
            }
        } else {
            var5 = -0.24278529
        }
    }
    var var6 float64
    if input[1] < 2.0 {
        if input[2] < 1.0 {
            if input[0] < 3.0 {
                if input[0] < 2.0 {
                    var6 = -0.06382415
                } else {
                    var6 = 0.074236676
                }
            } else {
                var6 = -0.30608842
            }
        } else {
            if input[3] < 1.0 {
                if input[0] < 2.0 {
                    var6 = 0.09499448
                } else {
                    var6 = 0.22013925
                }
            } else {
                var6 = -0.23133874
            }
        }
    } else {
        if input[0] < 1.0 {
            if input[1] < 4.0 {
                if input[2] < 1.0 {
                    var6 = -0.12510473
                } else {
                    var6 = 0.085284
                }
            } else {
                var6 = -0.23193018
            }
        } else {
            var6 = -0.23344965
        }
    }
    var var7 float64
    if input[1] < 2.0 {
        if input[2] < 1.0 {
            if input[3] < 2.0 {
                if input[1] < 1.0 {
                    var7 = -0.0040699397
                } else {
                    var7 = -0.060965866
                }
            } else {
                if input[3] < 6.0 {
                    var7 = -0.24518466
                } else {
                    var7 = 0.129678
                }
            }
        } else {
            if input[3] < 1.0 {
                if input[0] < 2.0 {
                    var7 = 0.08025375
                } else {
                    var7 = 0.2090828
                }
            } else {
                var7 = -0.21344204
            }
        }
    } else {
        if input[0] < 1.0 {
            if input[1] < 4.0 {
                if input[2] < 1.0 {
                    var7 = -0.10867168
                } else {
                    var7 = 0.07149207
                }
            } else {
                var7 = -0.22463627
            }
        } else {
            var7 = -0.22630154
        }
    }
    var var8 float64
    if input[1] < 2.0 {
        if input[2] < 1.0 {
            if input[0] < 3.0 {
                if input[0] < 2.0 {
                    var8 = -0.044570256
                } else {
                    var8 = 0.072612494
                }
            } else {
                var8 = -0.2730709
            }
        } else {
            if input[3] < 1.0 {
                if input[0] < 1.0 {
                    var8 = 0.05462835
                } else {
                    var8 = 0.1333761
                }
            } else {
                var8 = -0.19817665
            }
        }
    } else {
        if input[0] < 1.0 {
            if input[1] < 4.0 {
                if input[1] < 3.0 {
                    var8 = -0.11662438
                } else {
                    var8 = -0.03630318
                }
            } else {
                var8 = -0.2188653
            }
        } else {
            var8 = -0.22071591
        }
    }
    var var9 float64
    if input[1] < 2.0 {
        if input[3] < 2.0 {
            if input[2] < 1.0 {
                if input[1] < 1.0 {
                    var9 = 0.005153579
                } else {
                    var9 = -0.044111222
                }
            } else {
                if input[0] < 2.0 {
                    var9 = 0.050329536
                } else {
                    var9 = 0.19380306
                }
            }
        } else {
            if input[3] < 6.0 {
                var9 = -0.23245418
            } else {
                if input[0] < 1.0 {
                    var9 = 0.27938733
                } else {
                    var9 = -0.36091048
                }
            }
        }
    } else {
        if input[0] < 1.0 {
            if input[1] < 4.0 {
                if input[3] < 1.0 {
                    var9 = -0.10467117
                } else {
                    var9 = -0.022792349
                }
            } else {
                var9 = -0.21418849
            }
        } else {
            var9 = -0.21627223
        }
    }
    var var10 float64
    if input[1] < 2.0 {
        if input[3] < 2.0 {
            if input[0] < 2.0 {
                if input[1] < 1.0 {
                    var10 = 0.020716561
                } else {
                    var10 = -0.05095094
                }
            } else {
                if input[1] < 1.0 {
                    var10 = -0.3407025
                } else {
                    var10 = 0.13520713
                }
            }
        } else {
            if input[3] < 6.0 {
                var10 = -0.22412883
            } else {
                if input[0] < 1.0 {
                    var10 = 0.24495475
                } else {
                    var10 = -0.3034016
                }
            }
        }
    } else {
        if input[0] < 1.0 {
            if input[1] < 4.0 {
                if input[3] < 2.0 {
                    var10 = -0.05050053
                } else {
                    var10 = -0.232319
                }
            } else {
                var10 = -0.21030083
            }
        } else {
            var10 = -0.21267565
        }
    }
    var var11 float64
    if input[1] < 4.0 {
        if input[3] < 2.0 {
            if input[2] < 1.0 {
                if input[1] < 2.0 {
                    var11 = -0.017970001
                } else {
                    var11 = -0.07524441
                }
            } else {
                if input[1] < 3.0 {
                    var11 = 0.046946436
                } else {
                    var11 = 0.16423681
                }
            }
        } else {
            if input[3] < 6.0 {
                var11 = -0.22181216
            } else {
                if input[0] < 1.0 {
                    var11 = 0.21414696
                } else {
                    var11 = -0.2518617
                }
            }
        }
    } else {
        var11 = -0.20894623
    }
    var var12 float64
    if input[1] < 4.0 {
        if input[3] < 2.0 {
            if input[0] < 3.0 {
                if input[1] < 2.0 {
                    var12 = -0.007078701
                } else {
                    var12 = -0.053509977
                }
            } else {
                var12 = -0.22965859
            }
        } else {
            if input[3] < 6.0 {
                var12 = -0.21628737
            } else {
                if input[0] < 1.0 {
                    var12 = 0.19269772
                } else {
                    var12 = -0.23158506
                }
            }
        }
    } else {
        var12 = -0.2064226
    }
    var var13 float64
    if input[1] < 4.0 {
        if input[3] < 2.0 {
            if input[2] < 1.0 {
                if input[0] < 3.0 {
                    var13 = -0.022054456
                } else {
                    var13 = -0.21989587
                }
            } else {
                if input[1] < 3.0 {
                    var13 = 0.04115152
                } else {
                    var13 = 0.14927603
                }
            }
        } else {
            if input[3] < 6.0 {
                var13 = -0.211651
            } else {
                if input[3] < 7.0 {
                    var13 = 0.17399311
                } else {
                    var13 = -0.22727351
                }
            }
        }
    } else {
        var13 = -0.20422196
    }
    var var14 float64
    if input[1] < 4.0 {
        if input[3] < 2.0 {
            if input[1] < 1.0 {
                if input[0] < 2.0 {
                    var14 = 0.023755329
                } else {
                    var14 = -0.28980586
                }
            } else {
                if input[0] < 2.0 {
                    var14 = -0.0363174
                } else {
                    var14 = 0.070644
                }
            }
        } else {
            if input[3] < 6.0 {
                var14 = -0.20764653
            } else {
                if input[0] < 1.0 {
                    var14 = 0.16480525
                } else {
                    var14 = -0.2061611
                }
            }
        }
    } else {
        var14 = -0.20224023
    }
    var var15 float64
    var15 = sigmoid(var0 + var1 + var2 + var3 + var4 + var5 + var6 + var7 + var8 + var9 + var10 + var11 + var12 + var13 + var14)
    return []float64{1.0 - var15, var15}
}
func sigmoid(x float64) float64 {
    if (x < 0.0) {
        z := math.Exp(x)
        return z / (1.0 + z)
    }
    return 1.0 / (1.0 + math.Exp(-x))
}



func (m *Model) Score(features []float32) float32 {
	input := make([]float64, len(features))
	for i, v := range features {
		input[i] = float64(v)
	}

	rawScore := score(input)[0]
	probability := 1.0 / (1.0 + math.Exp(-rawScore))
	
	return float32(probability)
}
