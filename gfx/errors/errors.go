package errors

import "errors"

var (
	ErrorFeatureNotImplemented          = errors.New("feature not implemented, check your syntax")
	ErrorCanNotProceed                  = errors.New("can not proceed treatment")
	ErrorSizeDiffers                    = errors.New("sizes differs can not proceed treatment")
	ErrorCoordinatesNotFound            = errors.New("coordinates not found")
	ErrorUndefinedMode                  = errors.New("undefined mode")
	ErrorMissingNumberOfImageToGenerate = errors.New("iteration is not set, cannot define the number of images to generate")
	ErrorSizeMismatch                   = errors.New("error width and height mismatch cannot perform action")
	ErrorSizeOverflow                   = errors.New("size overflow the image size capacity")
	ErrorColorNotFound                  = errors.New("color not found in palette")
	ErrorNotYetImplemented              = errors.New("function is not yet implemented")
	ErrorModeNotFound                   = errors.New("mode not found or not implemented")
	ErrorBadSize                        = errors.New("width height does not correspond to data size")
	ErrorWidthSizeNotAccepted           = errors.New("width accepted  8 or 16 pixels")
	ErrorCustomDimensionMustBeSet       = errors.New("you must set custom width and height")
	ErrorCriteriaNotFound               = errors.New("criteria not found")
)
