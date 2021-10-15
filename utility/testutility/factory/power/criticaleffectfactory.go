package power

import "github.com/chadius/terosbattleserver/entity/power"

// CriticalEffectOptions is used to create healing effects.
type CriticalEffectOptions struct {
	damage                    int
	criticalHitThresholdBonus int
}

// CriticalEffectFactory creates a CriticalEffectOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func CriticalEffectFactory() *CriticalEffectOptions {
	return &CriticalEffectOptions{
		damage:                    0,
		criticalHitThresholdBonus: 0,
	}
}

// DealsDamage sets the amount of critical hit damage.
func (c *CriticalEffectOptions) DealsDamage(damage int) *CriticalEffectOptions {
	c.damage = damage
	return c
}

// CriticalHitThresholdBonus changes how far the attack has to succeed in order to critically hit.
func (c *CriticalEffectOptions) CriticalHitThresholdBonus(thresholdBonus int) *CriticalEffectOptions {
	c.criticalHitThresholdBonus = thresholdBonus
	return c
}

// Build uses the CriticalEffectOptions to create a CriticalEffect.
func (c *CriticalEffectOptions) Build() *power.CriticalEffect {
	newCriticalEffect := &power.CriticalEffect{
		Damage:                    c.damage,
		CriticalHitThresholdBonus: c.criticalHitThresholdBonus,
	}
	return newCriticalEffect
}
