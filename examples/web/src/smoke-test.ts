/**
 * Smoke Tests for Generated Code
 *
 * Quick tests to verify generated code compiles and basic functionality works.
 * Run with: npm run build (TypeScript compilation will catch any issues)
 */

import { runAllStabilityTests } from './generation-stability-tests';

// Run stability tests on module load
console.log('='.repeat(60));
console.log('GENERATION STABILITY SMOKE TEST');
console.log('='.repeat(60));

try {
  runAllStabilityTests();
  console.log('='.repeat(60));
  console.log('✅ SMOKE TEST PASSED - All generated code is stable');
  console.log('='.repeat(60));
} catch (error) {
  console.log('='.repeat(60));
  console.error('❌ SMOKE TEST FAILED - Breaking changes detected');
  console.error(error);
  console.log('='.repeat(60));
  throw error; // Re-throw to fail the build
}
