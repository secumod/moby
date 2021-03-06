package daemon

import (
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/image"
)

// TagImage creates the tag specified by newTag, pointing to the image named
// imageName (alternatively, imageName can also be an image ID).
func (daemon *Daemon) TagImage(imageName, repository, tag string) error {
	imageID, _, err := daemon.GetImageIDAndOS(imageName)
	if err != nil {
		return err
	}

	newTag, err := reference.ParseNormalizedNamed(repository)
	if err != nil {
		return err
	}
	if tag != "" {
		if newTag, err = reference.WithTag(reference.TrimNamed(newTag), tag); err != nil {
			return err
		}
	}

	return daemon.TagImageWithReference(imageID, newTag)
}

// TagImageWithReference adds the given reference to the image ID provided.
func (daemon *Daemon) TagImageWithReference(imageID image.ID, newTag reference.Named) error {
	if err := daemon.referenceStore.AddTag(newTag, imageID.Digest(), true); err != nil {
		return err
	}

	if err := daemon.imageStore.SetLastUpdated(imageID); err != nil {
		return err
	}
	daemon.LogImageEvent(imageID.String(), reference.FamiliarString(newTag), "tag")
	return nil
}
